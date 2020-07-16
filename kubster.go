package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/spf13/viper"
)

var Config struct {
	Bind       string
	ReadyDelay int
}

var live = true
var ready = false

func main() {
	viper.SetDefault("bind", ":8080")
	viper.SetDefault("readyDelay", 10)
	viper.SetEnvPrefix("kubster")
	viper.AutomaticEnv()
	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/live", handleLivenessProbe)
	mux.HandleFunc("/ready", handleReadinessProbe)
	mux.HandleFunc("/set", handleSet)
	mux.HandleFunc("/", handleRoot)

	go getReady(Config.ReadyDelay)
	log.Printf("Listening on %v, ready in %v seconds...\n", Config.Bind, Config.ReadyDelay)
	log.Fatal(http.ListenAndServe(Config.Bind, handlers.CombinedLoggingHandler(os.Stdout, mux)))
}

func getReady(delay int) {
	time.Sleep(time.Duration(delay) * time.Second)
	log.Println("Got ready!")
	ready = true
}

func handleRoot(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(response, kubeLogo)
}

func handleLivenessProbe(response http.ResponseWriter, request *http.Request) {
	if live {
		fmt.Fprintln(response, "Alive and kickin!")
	} else {
		http.Error(response, "I'm dead in the water!", http.StatusServiceUnavailable)
	}
}

func handleReadinessProbe(response http.ResponseWriter, request *http.Request) {
	if ready {
		fmt.Fprintln(response, "Ready player one!")
	} else {
		http.Error(response, "Not quite ready yet!", http.StatusServiceUnavailable)
	}
}

func handleSet(response http.ResponseWriter, request *http.Request) {
	err := setFlag(&live, request.URL.Query().Get("live"))
	if err != nil {
		http.Error(response, fmt.Sprintf("Can't parse liveness status: %v\n", err), http.StatusBadRequest)
	}

	err = setFlag(&ready, request.URL.Query().Get("ready"))
	if err != nil {
		http.Error(response, fmt.Sprintf("Can't parse readiness status: %v\n", err), http.StatusBadRequest)
	}
}

func setFlag(flag *bool, valStr string) error {
	if valStr == "" {
		return nil
	}
	val, err := strconv.ParseBool(valStr)
	if err == nil {
		*flag = val
	}
	return err
}

var dockerLogo = `                    ##        .            
              ## ## ##       ==            
           ## ## ## ##      ===            
       /""""""""""""""""\___/ ===        
  ~~~ {~~ ~~~~ ~~~ ~~~~ ~~ ~ /  ===- ~~~   
       \______ o          __/            
         \    \        __/             
          \____\______/                
 
          |          |
       __ |  __   __ | _  __   _
      /  \| /  \ /   |/  / _\ | 
      \__/| \__/ \__ |\_ \__  |`

var kubeLogo = `                                ///////                                
                            //// @@@@@&////                            
                        /////@@@@@@@@@@@@@%////@                       
                   &////@@@@@@@@@@@@@@@@@@@@@@%////&                   
               @////(@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@////&               
           /////@@@@@@@@@@@@@@@@@@@/%@@@@@@@@@@@@@@@@@@ ////           
        ////%@@@@@@@@@@@@@@@@@@@@@///@@@@@@@@@@@@@@@@@@@@@@////@       
      // @@@@@@@@@@@@@@@@@@@@@@@@@%//@@@@@@@@@@@@@@@@@@@@@@@@@(//      
     (/ @@@@@@@@@@@@@@@@@@@@@@@@@@@//@@@@@@@@@@@@@@@@@@@@@@@@@@#//     
     //@@@@@@@@@@@@@@@@@@@@@@/////////////@@@@@@@@@@@@@@@@@@@@@@//     
     //@@@@@@@@@///@@@@@@&///////////////////#@@@@@@(///@@@@@@@@(/&    
    //@@@@@@@@@@@@///@@//////@@@@@///(@@@@//////@@////%@@@@@@@@@@//    
    //@@@@@@@@@@@@@@//////%@@@@@@%///&@@@@@@@//////%@@@@@@@@@@@@@ /(   
   //@@@@@@@@@@@@@@@///////@@@@@@#////@@@@@@ //////@@@@@@@@@@@@@@@//   
   //@@@@@@@@@@@@@@////@/////(@@@@////@@@%//////(///@@@@@@@@@@@@@@//   
  //#@@@@@@@@@@@@@ ///@@@@///////////////////%@@@////@@@@@@@@@@@@@@//  
  //@@@@@@@@@@@@@@///(@@@@@@#//////////////@@@@@@@///#@@@@@@@@@@@@@//  
 //#@@@@@@@@@@@@@%///@@@@@@@@////&@@@%////@@@@@@@@////@@@@@@@@@@@@@@// 
 //@@@@@@@@@@@@@@@///@@(/////////@@@@@//////////@@////@@@@@@@@@@@@@@// 
&//@@@@@@@@@@@@@@////////////////// //////////////////%@@@@@@@@@@@@@#//
//@@@@@@@@@///////////#%@@@@@@%/////////(@@@@@@@#///////////%@@@@@@@@//
//@@@@@@@@@@@@@@@@@////@@@@@@@/////%/////@@@@@@@%///@@@@@@@@@@@@@@@@@//
//@@@@@@@@@@@@@@@@@@////@@@@@//// @@@&////@@@@@//// @@@@@@@@@@@@@@@@@//
 //(@@@@@@@@@@@@@@@@@(////@@@///&@@@@@%////@@&////@@@@@@@@@@@@@@@@@(// 
  .//%@@@@@@@@@@@@@@@@@/////////@@@@@@@(//////// @@@@@@@@@@@@@@@@@///  
    ///@@@@@@@@@@@@@@@@@@#///////// //////////@@@@@@@@@@@@@@@@@@///    
      //&@@@@@@@@@@@@@@@@@///@//////////// //@@@@@@@@@@@@@@@@@%//      
       ///@@@@@@@@@@@@@@@(//@@@@@@@@@@@@@@@@//@@@@@@@@@@@@@@@///       
         // @@@@@@@@@@@@///@@@@@@@@@@@@@@@@@(//%@@@@@@@@@@@%//         
           //%@@@@@@@@@@@%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@//*          
            ///@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ //            
              //&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%//              
               *//@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@///               
                 ///@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@(//                 `
