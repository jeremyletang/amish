package conf

type Github struct {
	Token        string   `json:"token"`
	Repositories []string `json:"repositories"`
	Users        []string `json:"users"`
}

type Mysql struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

type Gmail struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
	URL          string `json:"url"`
}

type Slack struct {
	Token    string   `json:"token"`
	Channels []string `json:"channels"`
}

type Conf struct {
	Github    Github   `json:"github"`
	Mysql     Mysql    `json:"mysql"`
	Gmail     Gmail    `json:"gmail"`
	Refresh   string   `json:"refresh"`
	Notify    string   `json:"notify"`
	Listeners []string `json:"listeners"`
	UseGmail  bool     `json:"use_gmail"`
	Slack     Slack    `json:"slack"`
}

const Seelog = `
<seelog>
  <outputs>
    <console formatid="colored"/>
  </outputs>
  <formats>
    <format id="colored"  format="%Date(2006 Jan 02/3:04:05.00 PM MST) (%File) [%EscM(36)%LEVEL%EscM(39)] %Msg%n%EscM(0)"/>
  </formats>
</seelog>`
