package project

import (
	"fmt"

	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/FriendsOfShopware/shopware-cli/shop"
)

var projectExtensionInstallCmd = &cobra.Command{
	Use:   "install [name]",
	Short: "Install a extension",
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

		activateAfterInstall, _ := cmd.PersistentFlags().GetBool("activate")

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

			if extension.InstalledAt != nil {
				log.Infof("Extension %s is already installed", arg)
				continue
			}

			if _, err := client.ExtensionManager.InstallExtension(adminSdk.NewApiContext(cmd.Context()), extension.Type, extension.Name); err != nil {
				failed = true

				log.Errorf("Installation of %s failed with error: %v", extension.Name, err)
			}

			if activateAfterInstall {
				if _, err := client.ExtensionManager.ActivateExtension(adminSdk.NewApiContext(cmd.Context()), extension.Type, extension.Name); err != nil {
					failed = true

					log.Errorf("Activation of %s failed with error: %v", extension.Name, err)
				} else {
					log.Infof("Activated %s", extension.Name)
				}
			}

			log.Infof("Installed %s", extension.Name)
		}

		if failed {
			return fmt.Errorf("install failed")
		}

		return nil
	},
}

func init() {
	projectExtensionCmd.AddCommand(projectExtensionInstallCmd)
	projectExtensionInstallCmd.PersistentFlags().Bool("activate", false, "Activate the extension")
}
