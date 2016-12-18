/*
 * +----------------------------------------------------------+
 * |                       DEAPEXER                           |
 * +----------------------------------------------------------+
 * | Written by Xan Manning, PyratLabs 2016. http://pyrat.io/ |
 * +----------------------------------------------------------+
 * |                                                          |
 * | Released under the MIT License. Please see LICENSE       |
 * +----------------------------------------------------------+
 * |                                                          |
 * | This application has a very specific use case for        |
 * | dealing with redirecting the apex record of a domain to  |
 * | a subdomain. The main reason for this is that a domain   |
 * | apex record must not be a CNAME, it must be an A record. |
 * |                                                          |
 * | A good example of when to use this is when working with  |
 * | AWS Elastic Load Balancers when you want to avoid using  |
 * | Route 53. This has been something that I have had to     |
 * | work around with customers. Often the solution is to 301 |
 * | redirect in Apache on one of the load balanced EC2       |
 * | instances. This causes disproportionate load on EC2s.    |
 * |                                                          |
 * | What DeApexer allows you to do is host one shared        |
 * | service that generates these redirects for any domain    |
 * | with a low cost to disk I/O. Configuration is stored in  |
 * | memory so this is fairly light on resources.             |
 * |                                                          |
 * +----------------------------------------------------------+
 */

package main

/*
 * Import our required Go modules
 */
import (
    "encoding/json"
    "flag"
    "io/ioutil"
    "log"
    "log/syslog"
    "net/http"
    "regexp"
    "strconv"
)

/*
 * Define our Configuration and UrlOverride data struct
 *
 * These will be populated from JSON data so we need to define in the struct
 * the key value relationship.
 */
type Configuration struct {
    ListenPort          int             `json:"listen_port"`
    DefaultSubdomain    string          `json:"default_subdomain"`
    Debug               bool            `json:"debug"`
    LogRedirect         bool            `json:"log_redirects"`
    Overrides           []UrlOverride   `json:"url_overrides"`
}

type UrlOverride struct {
    Hostname            string          `json:"hostname"`
    Subdomain           string          `json:"subdomain"`
    DefaultPath         string          `json:"default_path"`
    PathOverrides       []PathOverride  `json:"path_overrides"`
}

type PathOverride struct {
    Source              string          `json:"path_source"`
    Destination         string          `json:"path_destination"`
}

/*
 * Define our global configuration variable.
 */
var config Configuration

/*
 * Our init() function, used to build the configuration and set up syslog.
 */
func init() {
    // Do we have a specific location for our configuration file?
    configFile := flag.String("c",
        "/etc/deapexer/config.json",
        "Configuration File")

    // Are we wanting to avoid using Syslog?
    noLogs := flag.Bool("n", false, "Don't Log to Syslog")

    // Parse flag input
    flag.Parse()

    // Set up our syslog, throw error if there is a problem.
    if *noLogs == false {
        logfile, syslog_err := syslog.New(syslog.LOG_NOTICE, "deapexer")
        check(syslog_err)

        // Make log function write to log file.
        log.SetOutput(logfile)
    }

    // Read config data file, throw error if there is a problem.
    config_data, err := ioutil.ReadFile(*configFile)
    check(err)

    // Unmarshal the configuration file.
    unmarshal_err := json.Unmarshal(config_data, &config)
    check(unmarshal_err)

    // Debugging:
    //      - config file name
    //      - listening port
    //      - default subdomain.
    if config.Debug == true {
        log.Print("configFile: ", *configFile)
        log.Print("Loaded Configuration")
        log.Print("Listening Port: ", config.ListenPort)
        log.Print("Default Subdomain: ", config.DefaultSubdomain)
    }
}

/*
 * Define our function for checking errors and throwing panic.
 */
func check(e error) {
    if e != nil {
        panic(e)
    }
}


/*
 * Our deapex() function, generates the redirect
 */
func deapex(w http.ResponseWriter, r *http.Request) {
    // Are we logging redirects?
    if config.Debug == true || config.LogRedirect == true {
        log.Print("Request Host: ", r.Host)
        log.Print("Request URI: ", r.RequestURI)
    }

    // Obtain our host and URI from request data.
    requestHost := r.Host
    requestURI := r.RequestURI

    // Create the regexp to match the port number
    rpPort := regexp.MustCompile(":[0-9]+")

    // Strip the port number out of request host
    hostName := rpPort.ReplaceAllString(requestHost, "")

    // Variable to contain our subdomain redirect.
    var redirectSubdomain string

    // Subdomain Redirect Key
    var oKey int

    // Variable to contain our path redirect.
    var redirectPath string

    // Iterate through all the overrides we are storing in memory.
    for i := 0 ; i < len(config.Overrides) ; i++ {
        // Does the hostname value match override hostname?
        if config.Overrides[i].Hostname == hostName {
            // Set our redirect subdomain to this override.
            redirectSubdomain = config.Overrides[i].Subdomain
            oKey = i
            break
        }
    }

    // If there is no redirect subdomain in overrides, use the default.
    if redirectSubdomain == "" {
        redirectSubdomain = config.DefaultSubdomain + "." + hostName
    } else {
        // Iterate through the path overrides for this hostname.
        for i := 0 ; i < len(config.Overrides[oKey].PathOverrides) ; i++ {
            // Does the request URI value match override path?
            if config.Overrides[oKey].PathOverrides[i].Source == requestURI {
                // Set our URI to this path.
                redirectPath = config.Overrides[oKey].PathOverrides[i].
                    Destination
                break
            }
        }
    }

    // How are we going to redirect the path?
    if redirectPath == "" {
        if config.Overrides[oKey].DefaultPath != "" {
            requestURI = config.Overrides[oKey].DefaultPath
        }
    } else {
        requestURI = redirectPath
    }

    // Debug, where are we going?
    if config.Debug == true || config.LogRedirect == true {
        log.Print("New Host: ", redirectSubdomain)
        log.Print("New Request URI: ", requestURI)
    }

    // Build our forward address.
    forwardAddr := "http://" + redirectSubdomain + requestURI

    // Serve our 301 redirect.
    http.Redirect(w, r, forwardAddr, 301)
}

/*
 * Our main() function, sets up the http listener to handle deapex() function.
 */
func main() {
    // Define the function deapex() as our handler for.
    http.HandleFunc("/", deapex)

    // Listen and serve HTTP on our defined port.
    initErr := http.ListenAndServe(":" + strconv.Itoa(config.ListenPort), nil)

    // If there was an error setting up HTTP listener then log and throw panic.
    if initErr != nil {
        if config.Debug == true {
            log.Fatal("ListenAndServe: ", initErr)
        }
        panic(initErr)
    }
}
