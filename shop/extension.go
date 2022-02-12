package shop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

func (c Client) RefreshExtensions(ctx context.Context) error {
	req, err := c.newRequest(ctx, http.MethodPost, "/api/_action/extension/refresh", nil)
	if err != nil {
		return errors.Wrap(err, "RefreshExtensions")
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return errors.Wrap(err, "RefreshExtensions")
	}

	if err := resp.Body.Close(); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("RefreshExtensions: expected 204 from api, but got %d", resp.StatusCode)
	}

	return nil
}

func (c Client) GetAvailableExtensions(ctx context.Context) (ExtensionList, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/_action/extension/installed", nil)
	if err != nil {
		return nil, errors.Wrap(err, "GetAvailableExtensions")
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "GetAvailableExtensions")
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.Wrap(err, "GetAvailableExtensions")
	}

	var extensions ExtensionList
	if err := json.Unmarshal(content, &extensions); err != nil {
		return nil, err
	}

	return extensions, nil
}

type ExtensionList []*ExtensionDetail

func (l ExtensionList) GetByName(name string) *ExtensionDetail {
	for _, detail := range l {
		if detail.Name == name {
			return detail
		}
	}

	return nil
}

type ExtensionDetail struct {
	Extensions             []interface{} `json:"extensions"`
	Id                     interface{}   `json:"id"`
	LocalId                string        `json:"localId"`
	Name                   string        `json:"name"`
	Label                  string        `json:"label"`
	Description            string        `json:"description"`
	ShortDescription       interface{}   `json:"shortDescription"`
	ProducerName           string        `json:"producerName"`
	License                string        `json:"license"`
	Version                string        `json:"version"`
	LatestVersion          interface{}   `json:"latestVersion"`
	Languages              []interface{} `json:"languages"`
	Rating                 interface{}   `json:"rating"`
	NumberOfRatings        int           `json:"numberOfRatings"`
	Variants               []interface{} `json:"variants"`
	Faq                    []interface{} `json:"faq"`
	Binaries               []interface{} `json:"binaries"`
	Images                 []interface{} `json:"images"`
	Icon                   interface{}   `json:"icon"`
	IconRaw                *string       `json:"iconRaw"`
	Categories             []interface{} `json:"categories"`
	Permissions            []interface{} `json:"permissions"`
	Active                 bool          `json:"active"`
	Type                   string        `json:"type"`
	IsTheme                bool          `json:"isTheme"`
	Configurable           bool          `json:"configurable"`
	PrivacyPolicyExtension interface{}   `json:"privacyPolicyExtension"`
	StoreLicense           interface{}   `json:"storeLicense"`
	StoreExtension         interface{}   `json:"storeExtension"`
	InstalledAt            *struct {
		Date         string `json:"date"`
		TimezoneType int    `json:"timezone_type"`
		Timezone     string `json:"timezone"`
	} `json:"installedAt"`
	UpdatedAt    interface{}   `json:"updatedAt"`
	Notices      []interface{} `json:"notices"`
	Source       string        `json:"source"`
	UpdateSource string        `json:"updateSource"`
}

func (e ExtensionDetail) Status() string {
	var text string

	switch {
	case e.Source == "store":
		text = "can be downloaded from store"
	case e.Active:
		text = "installed, activated"
	case e.InstalledAt != nil:
		text = "installed, not activated"
	default:
		text = "not installed, not activated"
	}

	return text
}

func (c *Client) InstallExtension(ctx context.Context, extType, name string) error {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/_action/extension/install/%s/%s", extType, name), nil)

	if err != nil {
		return errors.Wrap(err, "InstallExtension")
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf("InstallExtension: got http code %d from api: %s", resp.StatusCode, string(content))
	}

	return nil
}

func (c *Client) UninstallExtension(ctx context.Context, extType, name string) error {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/_action/extension/uninstall/%s/%s", extType, name), nil)

	if err != nil {
		return errors.Wrap(err, "UninstallExtension")
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf("UninstallExtension: got http code %d from api: %s", resp.StatusCode, string(content))
	}

	return nil
}

func (c *Client) ActivateExtension(ctx context.Context, extType, name string) error {
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/_action/extension/activate/%s/%s", extType, name), nil)

	if err != nil {
		return errors.Wrap(err, "ActivateExtension")
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf("ActivateExtension: got http code %d from api: %s", resp.StatusCode, string(content))
	}

	return nil
}

func (c *Client) DeactivateExtension(ctx context.Context, extType, name string) error {
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/_action/extension/deactivate/%s/%s", extType, name), nil)

	if err != nil {
		return errors.Wrap(err, "DeactivateExtension")
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf("DeactivateExtension: got http code %d from api: %s", resp.StatusCode, string(content))
	}

	return nil
}

func (c *Client) RemoveExtension(ctx context.Context, extType, name string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/_action/extension/remove/%s/%s", extType, name), nil)

	if err != nil {
		return errors.Wrap(err, "RemoveExtension")
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf("RemoveExtension: got http code %d from api: %s", resp.StatusCode, string(content))
	}

	return nil
}

func (c *Client) UpdateExtension(ctx context.Context, extType, name string) error {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/_action/extension/update/%s/%s", extType, name), nil)

	if err != nil {
		return errors.Wrap(err, "UpdateExtension")
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf("UpdateExtension: got http code %d from api: %s", resp.StatusCode, string(content))
	}

	return nil
}

func (c *Client) DownloadExtension(ctx context.Context, name string) error {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/_action/extension/download/%s", name), nil)

	if err != nil {
		return errors.Wrap(err, "DownloadExtension")
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf("DownloadExtension: got http code %d from api: %s", resp.StatusCode, string(content))
	}

	return nil
}

func (c *Client) UploadExtension(ctx context.Context, extensionZip io.Reader) error {
	var buf bytes.Buffer
	parts := multipart.NewWriter(&buf)
	mimeHeader := textproto.MIMEHeader{}
	mimeHeader.Set("Content-Disposition", `form-data; name="file"; filename="extension.zip"`)
	mimeHeader.Set("Content-Type", "application/zip")

	part, err := parts.CreatePart(mimeHeader)
	if err != nil {
		return err
	}

	if _, err := io.Copy(part, extensionZip); err != nil {
		return err
	}
	if err := parts.Close(); err != nil {
		return err
	}

	var body io.Reader = &buf

	req, err := c.newRequest(ctx, http.MethodPost, "/api/_action/extension/upload", body)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", parts.FormDataContentType())

	var resp *http.Response

	if resp, err = c.httpClient.Do(req); err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errors.New("could not upload extension")
	}

	return nil
}