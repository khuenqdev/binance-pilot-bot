package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "github.com/khuenqdev/binance-pilot-bot/binance"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
    Use:   "test",
    Short: "Test binance API calls",
    Run:   runTest,
}

func init() {
    rootCmd.AddCommand(testCmd)

    testCmd.Flags().StringP("apiKey", "", "Cu6DzaLqpTHLpf2NfBW9TH5TiVHawviX5nm3rB13CAcMUtObO3LNbU7UuW8BaR58", "Binance API key")
    testCmd.Flags().StringP("apiSecret", "", "Eww4fIIrKLMa68LNzYE7gHcNVA1SxkprpKFnST5HGm4G2bqr6QHT4SVHBKgoXQ3N", "Binance API secret")

    _ = viper.BindPFlag("address", testCmd.Flags().Lookup("apiKey"))
    _ = viper.BindPFlag("address", testCmd.Flags().Lookup("apiSecret"))
}

func runTest(cmd *cobra.Command, args []string) {
    apiKey := viper.GetString("apiKey")
    apiSecret := viper.GetString("apiSecret")

    apiClient := binance.NewClient(binance.BaseUrlTestnet, apiKey, apiSecret)
    err := apiClient.MarketService.Ping()

    if nil != err {
        panic(err)
    }

    fmt.Println("Success!")
}
