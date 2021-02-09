/*
Copyright 2017, 2019 the Velero contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"io/ioutil"

	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"strings"

	"github.com/drajshek/k8s-ingress-mutator/pkg/ingress2httpproxy"
	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	core "k8s.io/api/networking/v1beta1"
)

//var log = logrus.New()

const pluginName = "kubectl-mutate-ingress2httpproxy"

func main() {

	var log = &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.WarnLevel,
	}

	// command line parameters
	var ingressFilePath string
	flag.StringVar(&ingressFilePath, "f", "", "the name of the input file with Ingress . The default is stdin.")
	var httpProxyFilePath string
	flag.StringVar(&httpProxyFilePath, "o", "", "the name of the output file of HTTPProxy. The default is stdout.")
	debug := flag.Bool("debug", false, "debug flag, it turns verbose on")
	jsonFormat := flag.Bool("json", false, "json parameter, it accepts json format as input. The output will be in json as well.")

	flag.Parse()

	var lvl string
	if *debug {
		lvl = "debug"
	} else {
		lvl, _ = os.LookupEnv("LOG_LEVEL")
	}

	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}

	if *debug {
		log.SetLevel(ll)
	}

	if ingressFilePath == "" {
		log.Debugf("[%s] Reading Ingress from stdin.", pluginName)
	} else {
		log.Debugf("[%s] Ingress input file is %s", pluginName, ingressFilePath)
	}

	if httpProxyFilePath == "" {
		log.Debugf("[%s] Writing HTTPProxy to stdout.", pluginName)
	} else {
		log.Debugf("[%s] HTTPProxy output file is %s", pluginName, httpProxyFilePath)
	}
	log.Debugf("[%s] json input format file is set to %#v", pluginName, *jsonFormat)

	var in io.Reader

	if ingressFilePath != "" {
		f, err := os.Open(ingressFilePath)
		if err != nil {
			log.Errorf("[%s] %s", pluginName, err)
			os.Exit(1)
		}
		defer f.Close()
		in = f
	} else {
		in = os.Stdin
	}

	buf := bufio.NewScanner(in)

	var item string

	for buf.Scan() {
		item += buf.Text()
		item += "\n"
	}

	item = strings.TrimSuffix(item, "\n")

	log.Tracef("[%s] Input format file content: %s", pluginName, item)

	if err := buf.Err(); err != nil {
		log.Errorf("[%s] error reading file %s", pluginName, err)
	}

	ingress := core.Ingress{}

	if *jsonFormat {
		log.Debugf("[%s] Processing json input.", pluginName)

		itemMarshal := []byte(item)
		err := json.Unmarshal(itemMarshal, &ingress)
		if err != nil {
			log.Errorf("[%s] %#v", pluginName, err)
			if jsonErr, ok := err.(*json.SyntaxError); ok {
				problemPart := itemMarshal[jsonErr.Offset-10 : jsonErr.Offset+10]
				err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
			}
			return
		}
	} else {
		log.Debugf("[%s] Processing yaml input.", pluginName)
		itemMarshal := []byte(item)
		err := yaml.Unmarshal(itemMarshal, ingress)
		if err != nil {
			log.Errorf("[%s] %#v", pluginName, err)
			return
		}
	}

	var domain string

	domain = ingress.Spec.TLS[0].Hosts[0]

	hp := ingress2httpproxy.NewMutator(pluginName, log, ingress, domain)
	output := hp.Mutate()

	var outputHTTPProxy []byte
	if *jsonFormat {
		log.Debugf("[%s] Writing json output.", pluginName)
		outputHTTPProxy, _ = json.MarshalIndent(&output.HTTPProxy, "", " ")
	} else {
		log.Debugf("[%s] Writing YAML output.", pluginName)
		fmt.Print("ff", &output.HTTPProxy)
		outputHTTPProxy, _ = yaml.Marshal(&output.HTTPProxy)
	}

	if httpProxyFilePath == "" {
		fmt.Println(string(outputHTTPProxy))
	} else {
		err = ioutil.WriteFile(httpProxyFilePath, outputHTTPProxy, 0644)
		if err != nil {
			log.Errorf("[%s] %s", pluginName, err)
			os.Exit(1)
		}
	}

}
