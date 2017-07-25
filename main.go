// bosh is an interactive shell for working with Bankrs OS
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/bankrs/bosgo"
)

type state struct {
	client    *bosgo.Client
	devEmail  string
	devClient *bosgo.DevClient

	applicationID string
	appClient     *bosgo.AppClient

	userName   string
	userClient *bosgo.UserClient
}

var session state
var addr = flag.String("a", "api.sandbox.bankrs.com", "address of api to connect to")

func main() {
	flag.Parse()

	session.client = bosgo.New(http.DefaultClient, *addr)
	shell := ishell.New()
	shell.SetPrompt("> ")

	shell.AddCmd(&ishell.Cmd{
		Name: "createdev",
		Help: "create a new developer account",
		Func: createDeveloper,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "login",
		Help: "login with an existing developer account",
		Func: loginDeveloper,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "changepassword",
		Help: "change password for the current developer account",
		Func: changePasswordDeveloper,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "logout",
		Help: "logout from a developer account",
		Func: logoutDeveloper,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "deletedeveloper",
		Help: "delete a developer account",
		Func: deleteDeveloper,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "lostpassword",
		Help: "send a lost password request",
		Func: lostPassword,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "resetpassword",
		Help: "reset a lost password",
		Func: resetPassword,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "profile",
		Help: "show the developer's profile",
		Func: profileDeveloper,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "setprofile",
		Help: "set the developer's profile",
		Func: setProfileDeveloper,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "createapp",
		Help: "create an application",
		Func: createApplication,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "listapps",
		Help: "list registered applications",
		Func: listApplications,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "updateapp",
		Help: "update an application",
		Func: updateApplication,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "deleteapp",
		Help: "delete an application",
		Func: deleteApplication,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "useapp",
		Help: "switch to using an application",
		Func: useApplication,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "stats",
		Help: "display stats for a developer",
		Func: stats,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "createuser",
		Help: "create a new user",
		Func: createUser,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "listusers",
		Help: "list users",
		Func: listUsers,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "loginuser",
		Help: "login as a user",
		Func: loginUser,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "logoutuser",
		Help: "logout from a user session",
		Func: logoutUser,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "deleteuser",
		Help: "delete the current user",
		Func: deleteUser,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "categories",
		Help: "list classification categories",
		Func: categories,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "searchproviders",
		Help: "search financial providers",
		Func: searchProviders,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "provider",
		Help: "lookup a single financial provider",
		Func: provider,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "accesses",
		Help: "list bank accesses for a user",
		Func: accesses,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "addaccess",
		Help: "add a bank accesses for a user",
		Func: addAccess,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "deleteaccess",
		Help: "delete a bank accesses",
		Func: deleteAccess,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getaccess",
		Help: "get details of a bank accesses",
		Func: getAccess,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "updateaccess",
		Help: "update challenge answers for a bank access",
		Func: updateAccess,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "refreshaccesses",
		Help: "refresh all bank accesses",
		Func: refreshAccesses,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "job",
		Help: "show the status of a job",
		Func: job,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "answer",
		Help: "provide a challenge answer for a job",
		Func: answer,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "canceljob",
		Help: "cancel a job",
		Func: cancelJob,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "accounts",
		Help: "list bank accounts for a user",
		Func: accounts,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getaccount",
		Help: "get details of a single account",
		Func: getAccount,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "transactions",
		Help: "list transactions for a user",
		Func: transactions,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "gettransaction",
		Help: "get details of a single transaction",
		Func: getTransaction,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "scheduledtransactions",
		Help: "list scheduled transactions for a user",
		Func: scheduledTransactions,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getscheduledtransaction",
		Help: "get details of a single scheduled transaction",
		Func: getScheduledTransaction,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "repeatedtransactions",
		Help: "list repeated transactions for a user",
		Func: repeatedTransactions,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getrepeatedtransaction",
		Help: "get details of a single repeated transaction",
		Func: getRepeatedTransaction,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "deleterecurringtransfer",
		Help: "delete a recurring transfer",
		Func: deleteRecurringTransfer,
	})

	// run shell
	shell.Run()
}

func createDeveloper(c *ishell.Context) {
	email, password, err := readCredentials("Email", c)
	if err != nil {
		c.Err(err)
		return
	}

	devClient, err := session.client.CreateDeveloper(email, password).Send()
	if err != nil {
		c.Err(err)
		return
	}

	session.devClient = devClient
}

func loginDeveloper(c *ishell.Context) {
	email, password, err := readCredentials("Email", c)
	if err != nil {
		c.Err(err)
		return
	}

	devClient, err := session.client.Login(email, password).Send()
	if err != nil {
		c.Err(err)
		return
	}
	session.devEmail = email
	session.devClient = devClient
	c.SetPrompt(email + "> ")
}

func lostPassword(c *ishell.Context) {
	email := readArg(0, "Email", c)

	err := session.client.LostPassword(email).Send()
	if err != nil {
		c.Err(err)
		return
	}
}

func resetPassword(c *ishell.Context) {
	password := readArgPassword(0, "Password", c)
	token := readArg(1, "Token", c)

	err := session.client.ResetPassword(password, token).Send()
	if err != nil {
		c.Err(err)
		return
	}
}

func logoutDeveloper(c *ishell.Context) {
	if session.devClient == nil {
		c.Err(fmt.Errorf("login to a developer account first"))
		return
	}

	err := session.devClient.Logout().Send()
	if err != nil {
		c.Err(err)
		return
	}
	session.devEmail = ""
	session.devClient = nil
	c.SetPrompt("> ")
}

func deleteDeveloper(c *ishell.Context) {
	if session.devClient == nil {
		c.Err(fmt.Errorf("login to a developer account first"))
		return
	}

	err := session.devClient.Delete().Send()
	if err != nil {
		c.Err(err)
		return
	}
	session.devEmail = ""
	session.devClient = nil
	c.SetPrompt("> ")
}

func profileDeveloper(c *ishell.Context) {
	if session.devClient == nil {
		c.Err(fmt.Errorf("login to a developer account first"))
		return
	}

	profile, err := session.devClient.Profile().Send()
	if err != nil {
		c.Err(err)
		return
	}
	c.Printf("Company: %s\n", profile.Company)
	c.Printf("Has production access: %v\n", profile.HasProductionAccess)
}

func setProfileDeveloper(c *ishell.Context) {
	if session.devClient == nil {
		c.Err(fmt.Errorf("login to a developer account first"))
		return
	}

	var profile bosgo.DeveloperProfile
	profile.Company = readArg(0, "Company name", c)
	profile.HasProductionAccess = readArgBool(1, "Has production access (y/n)", c)

	err := session.devClient.SetProfile(&profile).Send()
	if err != nil {
		c.Err(err)
		return
	}
}

func changePasswordDeveloper(c *ishell.Context) {
	if session.devClient == nil {
		c.Err(fmt.Errorf("login to a developer account first"))
		return
	}

	if len(c.Args) < 2 {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)
	}

	var oldpwd, newpwd string
	if len(c.Args) < 1 {
		c.Print("Old password: ")
		oldpwd = c.ReadLine()
	} else {
		oldpwd = c.Args[0]
	}

	if len(c.Args) < 2 {
		c.Print("New password: ")
		newpwd = c.ReadPassword()
	} else {
		newpwd = c.Args[1]
	}

	err := session.devClient.ChangePassword(oldpwd, newpwd).Send()
	if err != nil {
		c.Err(err)
		return
	}
}

func createApplication(c *ishell.Context) {
	if session.devClient == nil {
		c.Err(fmt.Errorf("login to a developer account first"))
		return
	}

	label, err := readOneArg("Label", c)
	if err != nil {
		c.Err(err)
		return
	}

	appID, err := session.devClient.Applications.Create(label).Send()
	if err != nil {
		c.Err(err)
		return
	}
	c.Println("application id", appID)
}

func listApplications(c *ishell.Context) {
	if session.devClient == nil {
		c.Err(fmt.Errorf("login to a developer account first"))
		return
	}

	list, err := session.devClient.Applications.List().Send()
	if err != nil {
		c.Err(err)
		return
	}

	for _, app := range list.Applications {
		c.Printf("%s (%s)\n", app.Label, app.ApplicationID)
	}
}

func updateApplication(c *ishell.Context) {
	if session.devClient == nil {
		c.Err(fmt.Errorf("login to a developer account first"))
		return
	}

	applicationID := readArg(0, "Application ID", c)
	label := readArg(1, "Label", c)

	err := session.devClient.Applications.Update(applicationID, label).Send()
	if err != nil {
		c.Err(err)
		return
	}
}

func deleteApplication(c *ishell.Context) {
	if session.devClient == nil {
		c.Err(fmt.Errorf("login to a developer account first"))
		return
	}

	applicationID := readArg(0, "Application ID", c)

	err := session.devClient.Applications.Delete(applicationID).Send()
	if err != nil {
		c.Err(err)
		return
	}
}

func useApplication(c *ishell.Context) {
	appID, err := readOneArg("Application ID", c)
	if err != nil {
		c.Err(err)
		return
	}

	session.appClient = session.client.WithApplicationID(appID)
	session.applicationID = appID
	c.SetPrompt(appID + "> ")
}

func listUsers(c *ishell.Context) {
	if session.appClient == nil {
		c.Err(fmt.Errorf("use an application id first"))
		return
	}

	list, err := session.appClient.Users.List().Send()
	if err != nil {
		c.Err(err)
		return
	}

	for _, user := range list.Users {
		c.Printf("* %s\n", user)
	}
}

func stats(c *ishell.Context) {
	if session.devClient == nil {
		c.Err(fmt.Errorf("login to a developer account first"))
		return
	}

	statType := readArg(0, "Type", c)
	var fromDate, toDate time.Time
	if len(c.Args) > 2 {
		var err error
		if fromDate, err = time.Parse("2006-01-02", c.Args[1]); err != nil {
			c.Err(fmt.Errorf("expected a date in yyyy-mm-dd format: %v", err))
			return
		}

		if toDate, err = time.Parse("2006-01-02", c.Args[2]); err != nil {
			c.Err(fmt.Errorf("expected a date in yyyy-mm-dd format: %v", err))
			return
		}
	}

	switch strings.ToLower(statType) {
	case "merchants":
		req := session.devClient.Stats.Merchants()
		if !fromDate.IsZero() && !toDate.IsZero() {
			req.FromDate(fromDate)
			req.ToDate(toDate)
		}
		stats, err := req.Send()
		if err != nil {
			c.Err(err)
			return
		}
		dumpJSON(c, stats)
	case "providers":
		req := session.devClient.Stats.Providers()
		if !fromDate.IsZero() && !toDate.IsZero() {
			req.FromDate(fromDate)
			req.ToDate(toDate)
		}
		stats, err := req.Send()
		if err != nil {
			c.Err(err)
			return
		}
		dumpJSON(c, stats)
	case "transfers":
		req := session.devClient.Stats.Transfers()
		if !fromDate.IsZero() && !toDate.IsZero() {
			req.FromDate(fromDate)
			req.ToDate(toDate)
		}
		stats, err := req.Send()
		if err != nil {
			c.Err(err)
			return
		}
		dumpJSON(c, stats)
	case "users":
		req := session.devClient.Stats.Users()
		if !fromDate.IsZero() && !toDate.IsZero() {
			req.FromDate(fromDate)
			req.ToDate(toDate)
		}
		stats, err := req.Send()
		if err != nil {
			c.Err(err)
			return
		}
		dumpJSON(c, stats)
	case "requests":
		req := session.devClient.Stats.Requests()
		if !fromDate.IsZero() && !toDate.IsZero() {
			req.FromDate(fromDate)
			req.ToDate(toDate)
		}
		stats, err := req.Send()
		if err != nil {
			c.Err(err)
			return
		}
		dumpJSON(c, stats)
	default:
		c.Err(fmt.Errorf("unknown stat type"))
	}
}

func createUser(c *ishell.Context) {
	if session.appClient == nil {
		c.Err(fmt.Errorf("use an application id first"))
		return
	}

	userName := readArg(0, "Name", c)
	password := readArgPassword(1, "Password", c)

	userClient, err := session.appClient.Users.Create(userName, password).Send()
	if err != nil {
		c.Err(err)
		return
	}

	session.userClient = userClient
	session.userName = userName
	c.SetPrompt(session.applicationID + "/" + session.userName + "> ")
}

func loginUser(c *ishell.Context) {
	if session.appClient == nil {
		c.Err(fmt.Errorf("use an application id first"))
		return
	}

	userName := readArg(0, "Name", c)
	password := readArgPassword(1, "Password", c)

	userClient, err := session.appClient.Users.Login(userName, password).Send()
	if err != nil {
		c.Err(err)
		return
	}

	session.userClient = userClient
	session.userName = userName
	c.SetPrompt(session.applicationID + "/" + session.userName + "> ")
}

func logoutUser(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("not logged in as a user"))
		return
	}
	err := session.userClient.Logout().Send()
	if err != nil {
		c.Err(err)
		return
	}

	session.userClient = nil
	session.userName = ""
	c.SetPrompt(session.applicationID + "> ")
}

func deleteUser(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("not logged in as a user"))
		return
	}
	err := session.userClient.Delete().Send()
	if err != nil {
		c.Err(err)
		return
	}

	session.userClient = nil
	session.userName = ""
	c.SetPrompt(session.applicationID + "> ")
}

func categories(c *ishell.Context) {
	if session.appClient == nil {
		c.Err(fmt.Errorf("use an application id first"))
		return
	}

	list, err := session.appClient.Categories.List().Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, list)
}

func searchProviders(c *ishell.Context) {
	if session.appClient == nil {
		c.Err(fmt.Errorf("use an application id first"))
		return
	}

	query := readArg(0, "Query", c)

	list, err := session.appClient.Providers.Search(query).Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, list)
}

func provider(c *ishell.Context) {
	if session.appClient == nil {
		c.Err(fmt.Errorf("use an application id first"))
		return
	}

	id := readArg(0, "Provider ID", c)

	list, err := session.appClient.Providers.Get(id).Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, list)
}

func accesses(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	list, err := session.userClient.Accesses.List().Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, list)
}

func addAccess(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	providerID := readArg(0, "Provider ID", c)
	answers := promptChallengeAnswers(c)

	req := session.userClient.Accesses.Add(providerID)
	for _, answer := range answers {
		req.ChallengeAnswer(answer)
	}

	job, err := req.Send()
	if err != nil {
		c.Err(err)
		return
	}

	c.Println("Job URI:", job.URI)
}

func deleteAccess(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	idstr := readArg(0, "Access ID", c)
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		c.Err(err)
		return
	}

	deleted, err := session.userClient.Accesses.Delete(id).Send()
	if err != nil {
		c.Err(err)
		return
	}

	c.Println("Deleted ID:", deleted)
}

func getAccess(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	idstr := readArg(0, "Access ID", c)
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		c.Err(err)
		return
	}

	access, err := session.userClient.Accesses.Get(id).Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, access)
}

func updateAccess(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	id := readArg(0, "Access ID", c)
	answers := promptChallengeAnswers(c)

	req := session.userClient.Accesses.Update(id)
	for _, answer := range answers {
		req.ChallengeAnswer(answer)
	}

	access, err := req.Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, access)
}

func refreshAccesses(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	jobs, err := session.userClient.Accesses.Refresh().Send()
	if err != nil {
		c.Err(err)
		return
	}

	c.Println("Job URIs:")
	for _, job := range jobs {
		c.Println(" * ", job.URI)
	}
}

func job(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}
	uri := readArg(0, "Job URI", c)

	status, err := session.userClient.Jobs.Get(uri).Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, status)
}

func answer(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}
	uri := readArg(0, "Job URI", c)
	answers := promptChallengeAnswers(c)

	req := session.userClient.Jobs.Answer(uri)
	for _, answer := range answers {
		req.ChallengeAnswer(answer)
	}

	err := req.Send()
	if err != nil {
		c.Err(err)
		return
	}
}

func cancelJob(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}
	uri := readArg(0, "Job URI", c)

	err := session.userClient.Jobs.Cancel(uri).Send()
	if err != nil {
		c.Err(err)
		return
	}
}

func accounts(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	list, err := session.userClient.Accounts.List().Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, list)
}

func getAccount(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	id := readArg(0, "Account ID", c)

	account, err := session.userClient.Accounts.Get(id).Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, account)
}

func transactions(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	list, err := session.userClient.Transactions.List().Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, list)
}

func getTransaction(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	id := readArg(0, "Account ID", c)

	tx, err := session.userClient.Transactions.Get(id).Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, tx)
}

func scheduledTransactions(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	list, err := session.userClient.ScheduledTransactions.List().Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, list)
}

func getScheduledTransaction(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	id := readArg(0, "Account ID", c)

	tx, err := session.userClient.ScheduledTransactions.Get(id).Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, tx)
}

func repeatedTransactions(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	list, err := session.userClient.RepeatedTransactions.List().Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, list)
}

func getRepeatedTransaction(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	id := readArg(0, "Account ID", c)

	tx, err := session.userClient.RepeatedTransactions.Get(id).Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, tx)
}

func deleteRecurringTransfer(c *ishell.Context) {
	if session.userClient == nil {
		c.Err(fmt.Errorf("login as a user first"))
		return
	}

	id := readArg(0, "ID", c)
	answers := promptChallengeAnswers(c)

	req := session.userClient.RecurringTransfers.Delete(id)
	for _, answer := range answers {
		req.ChallengeAnswer(answer)
	}

	tx, err := req.Send()
	if err != nil {
		c.Err(err)
		return
	}

	dumpJSON(c, tx)
}

func dumpJSON(c *ishell.Context, v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		c.Err(err)
		return
	}
	c.Println(string(data))
}

func readCredentials(userPrompt string, c *ishell.Context) (string, string, error) {
	if len(c.Args) < 2 {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)
	}

	var email, password string
	if len(c.Args) < 1 {
		c.Print(userPrompt + ": ")
		email = c.ReadLine()
	} else {
		email = c.Args[0]
	}

	if len(c.Args) < 2 {
		c.Print("Password: ")
		password = c.ReadPassword()
	} else {
		password = c.Args[1]
	}

	return email, password, nil
}

func readOneArg(prompt string, c *ishell.Context) (string, error) {
	if len(c.Args) < 1 {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)
	}

	var arg string
	if len(c.Args) < 1 {
		c.Print(prompt + ": ")
		arg = c.ReadLine()
	} else {
		arg = c.Args[0]
	}

	return arg, nil
}

func readArg(index int, prompt string, c *ishell.Context) string {
	if len(c.Args) < (index + 1) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)
	}

	var arg string
	if len(c.Args) < (index + 1) {
		c.Print(prompt + ": ")
		arg = c.ReadLine()
	} else {
		arg = c.Args[index]
	}

	return arg
}

func readArgPassword(index int, prompt string, c *ishell.Context) string {
	if len(c.Args) < (index + 1) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)
	}

	var arg string
	if len(c.Args) < (index + 1) {
		c.Print(prompt + ": ")
		arg = c.ReadPassword()
	} else {
		arg = c.Args[index]
	}

	return arg
}

func readArgBool(index int, prompt string, c *ishell.Context) bool {
	if len(c.Args) < (index + 1) {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true)
	}

	var arg string
	if len(c.Args) > index {
		arg = c.Args[index]
	}

	for {
		switch strings.ToLower(arg) {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			if v, err := strconv.ParseBool(arg); err == nil {
				return v
			}
			c.Print(prompt + ": ")
			arg = c.ReadLine()
		}
	}
}

func promptBool(c *ishell.Context, prompt string) bool {
	for {
		c.Print(prompt + ": ")
		arg := c.ReadLine()

		switch strings.ToLower(arg) {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			if v, err := strconv.ParseBool(arg); err == nil {
				return v
			}
		}
	}
}

func promptChallengeAnswers(c *ishell.Context) bosgo.ChallengeAnswerList {
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	var answers bosgo.ChallengeAnswerList
	for {
		var answer bosgo.ChallengeAnswer

		c.Print("Challenge ID (q to quit): ")
		answer.ID = c.ReadLine()
		if strings.ToLower(answer.ID) == "q" {
			return answers
		}

		c.Print("Value: ")
		answer.Value = c.ReadLine()
		answer.Store = promptBool(c, "Store (y/n)")

		answers = append(answers, answer)
	}
}
