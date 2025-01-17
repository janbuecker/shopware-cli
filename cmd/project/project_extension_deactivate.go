package project

import (
	"fmt"

	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/FriendsOfShopware/shopware-cli/shop"
)

var projectExtensionDeactivateCmd = &cobra.Command{
	Use:   "deactivate [name]",
	Short: "Deactivate a extension",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg *shop.Config
		var err error

		if cfg, err = shop.ReadConfig(projectConfigPath); err != nil {
			return err
		}

		client, err := shop.NewShopClient(cmd.Context(), cfg)
		if err != nil {
			return err
		}

		extensions, _, err := client.ExtensionManager.ListAvailableExtensions(adminSdk.NewApiContext(cmd.Context()))

		if err != nil {
			return err
		}

		failed := false

		for _, arg := range args {
			extension := extensions.GetByName(arg)

			if extension == nil {
				failed = true
				log.Errorf("Cannot find extension by name %s", arg)
				continue
			}

			if !extension.Active {
				log.Infof("Extension %s is already deactivated", arg)
				continue
			}

			if _, err := client.ExtensionManager.DeactivateExtension(adminSdk.NewApiContext(cmd.Context()), extension.Type, extension.Name); err != nil {
				failed = true

				log.Errorf("Deactivation of %s failed with error: %v", extension.Name, err)
			}

			log.Infof("Deactivated %s", extension.Name)
		}

		if failed {
			return fmt.Errorf("deactivation failed")
		}

		return nil
	},
}

func init() {
	projectExtensionCmd.AddCommand(projectExtensionDeactivateCmd)
}
