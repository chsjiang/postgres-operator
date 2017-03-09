/*
 Copyright 2017 Crunchy Data Solutions, Inc.
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
package cmd

import (
	"fmt"
	"github.com/crunchydata/operator/tpr"
	"github.com/spf13/cobra"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
)

// ShowCmd represents the show command
var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show a description of a database or cluster",
	Long: `show allows you to show the details of a database or cluster.
For example:

crunchy show database mydatabase
crunchy show cluster mycluster`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(`You must specify the type of resource to show.  Valid resource types include:
			        * database
				* cluster`)
		} else {
			switch args[0] {
			case "database":
			case "cluster":
				break
			default:
				fmt.Println(`You must specify the type of resource to show.  Valid resource types include:
			        * database
				* cluster`)
			}
		}

	},
}

func init() {
	fmt.Println("show init called")
	RootCmd.AddCommand(ShowCmd)
	ShowCmd.AddCommand(ShowDatabaseCmd)
	ShowCmd.AddCommand(ShowClusterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ShowCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ShowCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// showDatbaseCmd represents the show database command
var ShowDatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Show database information",
	Long: `Show a crunchy database. For example:

				crunchy show database mydatabase`,
	Run: func(cmd *cobra.Command, args []string) {
		showDatabase(args)
	},
}

// ShowClusterCmd represents the show cluster command
var ShowClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Show cluster information",
	Long: `Show a crunchy cluster. For example:

				crunchy show cluster mycluster`,
	Run: func(cmd *cobra.Command, args []string) {
		showCluster(args)
	},
}

func showDatabase(args []string) {
	//get a list of all databases
	databaseList := tpr.CrunchyDatabaseList{}
	err := Tprclient.Get().Resource("crunchydatabases").Do().Into(&databaseList)
	if err != nil {
		fmt.Println("error getting list of databases")
		fmt.Println(err.Error())
		return
	}

	//each arg represents a database name or the special 'all' value
	var pod *v1.Pod
	var service *v1.Service
	for _, arg := range args {
		fmt.Println("show database " + arg)
		for _, database := range databaseList.Items {
			if arg == "all" || database.Spec.Name == arg {
				fmt.Println("database LIST: " + database.Spec.Name)
				pod, err = Clientset.Core().Pods(api.NamespaceDefault).Get(database.Spec.Name)
				if err != nil {
					fmt.Println("error in getting database pod " + database.Spec.Name)
					fmt.Println(err.Error())
				} else {
					fmt.Println("pod " + pod.Name)
				}

				service, err = Clientset.Core().Services(api.NamespaceDefault).Get(database.Spec.Name)
				if err != nil {
					fmt.Println("error in getting database service " + database.Spec.Name)
					fmt.Println(err.Error())
				} else {
					fmt.Println("service " + service.Name)
				}
			}
		}
	}
}

func showCluster(args []string) {
	//get a list of all clusters
	clusterList := tpr.CrunchyClusterList{}
	err := Tprclient.Get().Resource("crunchyclusters").Do().Into(&clusterList)
	if err != nil {
		fmt.Println("error getting list of clusters")
		fmt.Println(err.Error())
		return
	}

	//each arg represents a cluster name or the special 'all' value
	var pod *v1.Pod
	var service *v1.Service
	for _, arg := range args {
		fmt.Println("show cluster " + arg)
		for _, cluster := range clusterList.Items {
			if arg == "all" || cluster.Spec.Name == arg {
				fmt.Println("cluster LIST: " + cluster.Spec.Name)
				pod, err = Clientset.Core().Pods(api.NamespaceDefault).Get(cluster.Spec.Name)
				if err != nil {
					fmt.Println("error in getting cluster pod " + cluster.Spec.Name)
					fmt.Println(err.Error())
				} else {
					fmt.Println("pod " + pod.Name)
				}

				service, err = Clientset.Core().Services(api.NamespaceDefault).Get(cluster.Spec.Name)
				if err != nil {
					fmt.Println("error in getting cluster service " + cluster.Spec.Name)
					fmt.Println(err.Error())
				} else {
					fmt.Println("service " + service.Name)
				}
			}
		}
	}

}

func ListPods() {
	//ConnectToKube()

	lo := v1.ListOptions{LabelSelector: "k8s-app=kube-dns"}
	fmt.Println("label selector is " + lo.LabelSelector)
	pods, err := Clientset.Core().Pods("").List(lo)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	for _, pod := range pods.Items {
		fmt.Println("pod Name " + pod.ObjectMeta.Name)
		fmt.Println("pod phase is " + pod.Status.Phase)
	}

}
