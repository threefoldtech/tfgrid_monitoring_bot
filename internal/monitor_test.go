package internal

import (
	"os"
	"testing"

	client "github.com/threefoldtech/substrate-client"
)

func TestParsers(t *testing.T) {
	t.Run("test_no_file", func(t *testing.T) {
		_, err := readFile("env.env")

		if err == nil {
			t.Errorf("expected error reading env.env")
		}
	})

	t.Run("test_valid_file", func(t *testing.T) {
		_, err := readFile("monitor.go")

		if err != nil {
			t.Errorf("expected no error, %v", err)
		}
	})

	t.Run("test_wrong_env_no_test_mnemonic", func(t *testing.T) {
		envContent := `
			TESTNET_MNEMONIC=
			MAINNET_MNEMONIC=mnemonic
			TFTS_LIMIT=700
			BOT_TOKEN=TOKEN
			CHAT_ID=ID
			MINS=1
		`

		_, err := parseEnv(envContent)

		if err == nil {
			t.Errorf("expected error, no TESTNET_MNEMONIC")
		}
	})

	t.Run("test_wrong_env_no_main_mnemonic", func(t *testing.T) {
		envContent := `
			TESTNET_MNEMONIC=mnemonic
			MAINNET_MNEMONIC=
			TFTS_LIMIT=700
			BOT_TOKEN=TOKEN
			CHAT_ID=ID
			MINS=1
		`

		_, err := parseEnv(envContent)

		if err == nil {
			t.Errorf("expected error, no MAINNET_MNEMONIC")
		}
	})

	t.Run("test_wrong_env_0_limit", func(t *testing.T) {
		envContent := `
			TESTNET_MNEMONIC=mnemonic
			MAINNET_MNEMONIC=mnemonic
			TFTS_LIMIT=0
			BOT_TOKEN=token
			CHAT_ID=ID
			MINS=1
		`

		_, err := parseEnv(envContent)

		if err == nil {
			t.Errorf("expected error, TFTS_LIMIT is 0")
		}
	})

	t.Run("test_wrong_env_no_token", func(t *testing.T) {
		envContent := `
			TESTNET_MNEMONIC=mnemonic
			MAINNET_MNEMONIC=mnemonic
			TFTS_LIMIT=8
			BOT_TOKEN=
			CHAT_ID=ID
			MINS=1
		`

		_, err := parseEnv(envContent)

		if err == nil {
			t.Errorf("expected error, BOT_TOKEN is missing")
		}
	})

	t.Run("test_wrong_env_no_id", func(t *testing.T) {
		envContent := `
			TESTNET_MNEMONIC=mnemonic
			MAINNET_MNEMONIC=mnemonic
			TFTS_LIMIT=8
			BOT_TOKEN=token
			CHAT_ID=
			MINS=1
		`

		_, err := parseEnv(envContent)

		if err == nil {
			t.Errorf("expected error, CHAT_ID is missing")
		}
	})

	t.Run("test_wrong_env_0_mins", func(t *testing.T) {
		envContent := `
			TESTNET_MNEMONIC=mnemonic
			MAINNET_MNEMONIC=mnemonic
			TFTS_LIMIT=8
			BOT_TOKEN=token
			CHAT_ID=id
			MINS=0
		`

		_, err := parseEnv(envContent)

		if err == nil {
			t.Errorf("expected error, MINS is 0")
		}
	})

	t.Run("test_wrong_env_string_mins", func(t *testing.T) {
		envContent := `
			TESTNET_MNEMONIC=mnemonic
			MAINNET_MNEMONIC=mnemonic
			TFTS_LIMIT=8
			BOT_TOKEN=token
			CHAT_ID=id
			MINS=min
		`

		_, err := parseEnv(envContent)

		if err == nil {
			t.Errorf("expected error, MINS is string")
		}
	})

	t.Run("test_wrong_env_string_limit", func(t *testing.T) {
		envContent := `
			TESTNET_MNEMONIC=mnemonic
			MAINNET_MNEMONIC=mnemonic
			TFTS_LIMIT=limit
			BOT_TOKEN=token
			CHAT_ID=id
			MINS=10
		`
		_, err := parseEnv(envContent)

		if err == nil {
			t.Errorf("expected error, TFTS_LIMIT is string")
		}
	})

	t.Run("test_wrong_env_key", func(t *testing.T) {
		envContent := `
			key=key
			TESTNET_MNEMONIC=mnemonic
			MAINNET_MNEMONIC=mnemonic
			TFTS_LIMIT=10
			BOT_TOKEN=token
			CHAT_ID=id
			MINS=10
		`
		_, err := parseEnv(envContent)

		if err == nil {
			t.Errorf("expected error, key is invalid")
		}
	})

	t.Run("test_valid_env", func(t *testing.T) {
		envContent := `
			TESTNET_MNEMONIC=mnemonic
			MAINNET_MNEMONIC=mnemonic
			TFTS_LIMIT=10
			BOT_TOKEN=token
			CHAT_ID=id
			MINS=10
		`
		_, err := parseEnv(envContent)

		if err != nil {
			t.Errorf("parsing should be successful")
		}
	})

	t.Run("test_valid_json", func(t *testing.T) {
		content := `
		{ 
			"mainnet": [],
			"testnet": []   
		}
		`
		_, err := parseJsonIntoWallets([]byte(content))

		if err != nil {
			t.Errorf("parsing should be successful")
		}
	})

	t.Run("test_invalid_json_no_test", func(t *testing.T) {
		content := `
		{ 
			"mainnet": []
		}
		`
		_, err := parseJsonIntoWallets([]byte(content))

		if err == nil {
			t.Errorf("parsing should fail, missing testnet")
		}
	})

	t.Run("test_invalid_json_no_main", func(t *testing.T) {
		content := `
		{ 
			"testnet": []
		}
		`
		_, err := parseJsonIntoWallets([]byte(content))

		if err == nil {
			t.Errorf("parsing should fail, missing mainnet")
		}
	})
}

func TestMonitor(t *testing.T) {
	//json
	jsonFile, err := os.CreateTemp("", "*.json")

	if err != nil {
		t.Errorf("failed with error, %v", err)
	}

	defer jsonFile.Close()
	defer os.Remove(jsonFile.Name())

	data := []byte(`{ 
		"mainnet": [""],
		"testnet": [""]   
	}`)
	if _, err := jsonFile.Write(data); err != nil {
		t.Error(err)
	}

	//env
	envFile, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Errorf("failed with error, %v", err)
	}

	defer envFile.Close()
	defer os.Remove(envFile.Name())

	data = []byte(`TESTNET_MNEMONIC=mnemonic
	MAINNET_MNEMONIC=mnemonic
	TFTS_LIMIT=1
	BOT_TOKEN=token
	CHAT_ID=id
	MINS=10`)
	if _, err := envFile.Write(data); err != nil {
		t.Error(err)
	}

	//managers
	substrate := map[network]client.Manager{}

	substrate[mainNetwork] = client.NewManager(SUBSTRATE_URLS[mainNetwork]...)
	substrate[testNetwork] = client.NewManager(SUBSTRATE_URLS[testNetwork]...)

	t.Run("test_invalid_monitor_env", func(t *testing.T) {
		_, err := NewMonitor("env", jsonFile.Name())

		if err == nil {
			t.Errorf("monitor should fail, wrong env")
		}
	})

	t.Run("test_invalid_monitor_json", func(t *testing.T) {

		_, err := NewMonitor(envFile.Name(), "wallets")

		if err == nil {
			t.Errorf("monitor should fail, wrong json")
		}
	})

	t.Run("test_valid_monitor", func(t *testing.T) {

		_, err := NewMonitor(envFile.Name(), jsonFile.Name())

		if err != nil {
			t.Errorf("monitor should be successful")
		}
	})

	t.Run("test_invalid_monitor_token", func(t *testing.T) {

		monitor, err := NewMonitor(envFile.Name(), jsonFile.Name())

		if err != nil {
			t.Errorf("monitor should be successful")
		}

		monitor.env.botToken = ""
		monitor.env.tftLimit = 10000000000000
		err = monitor.sendMessage(substrate[testNetwork], "")
		if err == nil {
			t.Errorf("sending a message should fail")
		}
	})

	t.Run("test_send_message_low_limit", func(t *testing.T) {

		monitor, err := NewMonitor(envFile.Name(), jsonFile.Name())

		if err != nil {
			t.Errorf("monitor should be successful")
		}

		err = monitor.sendMessage(substrate[testNetwork], "")
		if err == nil {
			t.Errorf("no message should be sent")
		}
	})
}

func TestWrongFilesContent(t *testing.T) {
	//json
	jsonFile, err := os.CreateTemp("", "*.json")

	if err != nil {
		t.Errorf("failed with error, %v", err)
	}

	defer jsonFile.Close()
	defer os.Remove(jsonFile.Name())

	data := []byte(`{ 
		"mainnet": []  
	}`)
	if _, err := jsonFile.Write(data); err != nil {
		t.Error(err)
	}

	t.Run("test_invalid_monitor_wrong_env", func(t *testing.T) {
		//env
		envFile, err := os.CreateTemp("", "*.env")
		if err != nil {
			t.Errorf("failed with error, %v", err)
		}

		defer envFile.Close()
		defer os.Remove(envFile.Name())

		data = []byte(`TESTNET_MNEMONIC=mnemonic`)
		if _, err := envFile.Write(data); err != nil {
			t.Error(err)
		}

		_, err = NewMonitor(envFile.Name(), jsonFile.Name())

		if err == nil {
			t.Errorf("monitor should fail, wrong env")
		}
	})

	t.Run("test_invalid_monitor_wrong_json", func(t *testing.T) {

		//env
		envFile, err := os.CreateTemp("", "*.env")
		if err != nil {
			t.Errorf("failed with error, %v", err)
		}

		defer envFile.Close()
		defer os.Remove(envFile.Name())

		data = []byte(`TESTNET_MNEMONIC=mnemonic
			MAINNET_MNEMONIC=mnemonic
			TFTS_LIMIT=1
			BOT_TOKEN=token
			CHAT_ID=id
			MINS=10`)
		if _, err := envFile.Write(data); err != nil {
			t.Error(err)
		}

		_, err = NewMonitor(envFile.Name(), jsonFile.Name())

		if err == nil {
			t.Errorf("monitor should fail, wrong json")
		}
	})
}
