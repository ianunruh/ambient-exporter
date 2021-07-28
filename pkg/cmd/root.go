package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ianunruh/ambient-exporter/pkg/ambient"
	"github.com/ianunruh/ambient-exporter/pkg/collect"
)

var listenAddr string

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		log, err := zap.NewProduction()
		if err != nil {
			return err
		}

		apiKey := os.Getenv("AMBIENT_API_KEY")
		if apiKey == "" {
			return errors.New("AMBIENT_API_KEY env var is required")
		}

		appKey := os.Getenv("AMBIENT_APP_KEY")
		if appKey == "" {
			return errors.New("AMBIENT_APP_KEY env var is required")
		}

		httpClient := &http.Client{
			Timeout: 10 * time.Second,
		}

		client := ambient.NewClient(apiKey, appKey, httpClient)

		prometheus.MustRegister(collect.NewCollector(client, log))

		http.Handle("/metrics", promhttp.Handler())

		log.Info("Starting metrics server",
			zap.String("address", listenAddr))
		return http.ListenAndServe(listenAddr, nil)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&listenAddr, "listen", "l", ":9090", "Host/port to listen on")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
