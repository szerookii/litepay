package litecoin

type BlockchainInfo struct {
	Chain                string  `json:"chain"`
	Blocks               int     `json:"blocks"`
	Headers              int     `json:"headers"`
	Bestblockhash        string  `json:"bestblockhash"`
	Difficulty           float64 `json:"difficulty"`
	Mediantime           int     `json:"mediantime"`
	Verificationprogress float64 `json:"verificationprogress"`
	Initialblockdownload bool    `json:"initialblockdownload"`
	Chainwork            string  `json:"chainwork"`
	SizeOnDisk           int64   `json:"size_on_disk"`
	Pruned               bool    `json:"pruned"`
	Softforks            struct {
		Bip34 struct {
			Type   string `json:"type"`
			Active bool   `json:"active"`
			Height int    `json:"height"`
		} `json:"bip34"`
		Bip66 struct {
			Type   string `json:"type"`
			Active bool   `json:"active"`
			Height int    `json:"height"`
		} `json:"bip66"`
		Bip65 struct {
			Type   string `json:"type"`
			Active bool   `json:"active"`
			Height int    `json:"height"`
		} `json:"bip65"`
		Csv struct {
			Type   string `json:"type"`
			Active bool   `json:"active"`
			Height int    `json:"height"`
		} `json:"csv"`
		Segwit struct {
			Type   string `json:"type"`
			Active bool   `json:"active"`
			Height int    `json:"height"`
		} `json:"segwit"`
		Taproot struct {
			Type string `json:"type"`
			Bip8 struct {
				Status        string `json:"status"`
				StartHeight   int    `json:"start_height"`
				TimeoutHeight int    `json:"timeout_height"`
				Since         int    `json:"since"`
			} `json:"bip8"`
			Height int  `json:"height"`
			Active bool `json:"active"`
		} `json:"taproot"`
		Mweb struct {
			Type string `json:"type"`
			Bip8 struct {
				Status        string `json:"status"`
				StartHeight   int    `json:"start_height"`
				TimeoutHeight int    `json:"timeout_height"`
				Since         int    `json:"since"`
			} `json:"bip8"`
			Height int  `json:"height"`
			Active bool `json:"active"`
		} `json:"mweb"`
	} `json:"softforks"`
	Warnings string `json:"warnings"`
}

type AddressInfo struct {
	Address             string   `json:"address"`
	Label               string   `json:"label"`
	ScriptPubKey        string   `json:"scriptPubKey"`
	Ismine              bool     `json:"ismine"`
	Solvable            bool     `json:"solvable"`
	Desc                string   `json:"desc"`
	Iswatchonly         bool     `json:"iswatchonly"`
	Isscript            bool     `json:"isscript"`
	Iswitness           bool     `json:"iswitness"`
	WitnessVersion      int      `json:"witness_version"`
	WitnessProgram      string   `json:"witness_program"`
	Ismweb              bool     `json:"ismweb"`
	Pubkey              string   `json:"pubkey"`
	Ischange            bool     `json:"ischange"`
	Timestamp           int      `json:"timestamp"`
	Hdkeypath           string   `json:"hdkeypath"`
	Hdseedid            string   `json:"hdseedid"`
	Hdmasterfingerprint string   `json:"hdmasterfingerprint"`
	Labels              []string `json:"labels"`
}
