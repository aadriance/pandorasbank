package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type CombatPos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type StatStages struct {
	PhysicalAttack  int `json:"physicalAttack"`
	AstralAttack    int `json:"astralAttack"`
	PhysicalDefense int `json:"physicalDefense"`
	AstralDefense   int `json:"astralDefense"`
	Speed           int `json:"speed"`
	Accuracy        int `json:"accuracy"`
	Evasion         int `json:"evasion"`
}

type Elestral struct {
	ID struct {
		SerializedVersion string `json:"serializedVersion"`
		Hash              string `json:"Hash"`
	} `json:"id"`
	Name                        string     `json:"name"`
	Species                     string     `json:"species"`
	UsesStellarMaterial         bool       `json:"usesStellarMaterial"`
	Element                     int        `json:"element"`
	SubElement                  int        `json:"subElement"`
	HealthBaseStat              int        `json:"healthBaseStat"`
	PhysicalAttack              int        `json:"physicalAttack"`
	SpecialAttack               int        `json:"specialAttack"`
	PhysicalDefense             int        `json:"physicalDefense"`
	SpecialDefense              int        `json:"specialDefense"`
	Speed                       int        `json:"speed"`
	Ability0Name                string     `json:"Ability0Name"`
	Ability1Name                string     `json:"Ability1Name"`
	Ability2Name                string     `json:"Ability2Name"`
	Ability3Name                string     `json:"Ability3Name"`
	EmpoweredAbilityName        string     `json:"empoweredAbilityName"`
	CurrentLevel                int        `json:"currentLevel"`
	Health                      int        `json:"health"`
	MaxHealth                   int        `json:"maxHealth"`
	IsActiveCombat              bool       `json:"isActiveInCombat"`
	IsDodgeEnabled              bool       `json:"isDodgeEnabled"`
	IsCaster                    bool       `json:"isCaster"`
	IsStellar                   bool       `json:"isStellar"`
	TeamSlot                    int        `json:"teamSlot"`
	AbilityHitIndex             int        `json:"abilityHitIndex"`
	BondMeter                   int        `json:"bondMeter"`
	TurnOrderUIIndex            int        `json:"turnOrderUIIndex"`
	SelectedAbilityIndex        int        `json:"selectedAbilityIndex"`
	ShouldSkipTurn              bool       `json:"shouldSkipTurn"`
	LastDodgeTime               float64    `json:"lastDodgeTime"`
	LastSuccessfulDodgeTime     float64    `json:"lastSuccessfulDodgeTime"`
	CombatPos                   CombatPos  `json:"CombatPos"`
	CharacterType               int        `json:"characterType"`
	StatStages                  StatStages `json:"statStages"`
	SuppressedSlots             int        `json:"suppressedSlots"`
	SlotsUsedThisBattle         int        `json:"slotsUsedThisBattle"`
	IncomingDamageMultiplier    float64    `json:"IncomingDamageMultiplier"`
	OutgoingDamageMultiplier    float64    `json:"OutgoingDamageMultiplier"`
	MovesPerformedSinceLastSwap int        `json:"MovesPerformedSinceLastSwap"`
	LastUsedAbilitySlot         int        `json:"lastUsedAbilitySlot"`
	HasUsedEmpoweredAbility     bool       `json:"hasUsedEmpoweredAbility"`
}

type PlayerData struct {
	Name            string    `json:"Name"`
	SpiritElement   int       `json:"SpiritElement"`
	IsMaleCharacter bool      `json:"isMaleCharacter"`
	Money           int       `json:"Money"`
	FocusedSlot     int       `json:"FocusedSlot"`
	Character0      *Elestral `json:"Character0"`
	Character1      *Elestral `json:"Character1"`
	Character2      *Elestral `json:"Character2"`
	Character3      *Elestral `json:"Character3"`
	MaxSp           int       `json:"MaxCasterSP"`
	CurrentSp       int       `json:"CasterSP"`
	BondMeter       int       `json:"BondMeter"`
}

type StorageEntry struct {
	CharacterData *Elestral `json:"CharacterData"`
}

type StorageBox struct {
	Entries []StorageEntry `json:"entries"`
}

type ActiveBoons struct {
	ActiveBoonNames []string `json:"ActiveBoonNames"`
	BoonUsageCounts []int    `json:"BoonUsageCounts"`
}

type GameSave struct {
	ActivePlayerData PlayerData   `json:"activePlayerData"`
	StorageBoxes     []StorageBox `json:"storageBoxes"`
	GameFlags        []string     `json:"gameFlags"`
	ActiveBoons      ActiveBoons  `json:"activeBoons"`
	CurrentSceneName string       `json:"currentSceneName"`
	SaveVersion      string       `json:"saveVersion"`
	SaveTimestamp    string       `json:"saveTimestamp"`
}

type Settings struct {
	CustomSavePath string `json:"customSavePath"`
}

type Bank struct {
	Elestrals []*Elestral `json:"elestrals"`
}

func getElementName(element int) string {
	elements := map[int]string{
		0: "N/A",
		1: "Earth",
		2: "Fire",
		3: "Water",
		4: "Thunder",
		5: "Wind",
		6: "Frost",
		7: "Solar",
		8: "Lunar",
	}
	if name, ok := elements[element]; ok {
		return name
	}
	return fmt.Sprintf("Unknown (%d)", element)
}

func createElestralCard(e *Elestral, onSave func(), onExport func(), onImport func(), onRelease func()) *widget.Card {
	if e == nil || e.Species == "" {
		return nil
	}

	stellar := ""
	if e.IsStellar {
		stellar = " (Stellar)"
	}

	nameLabel := widget.NewLabel(e.Name)
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}
	editButton := widget.NewButton("Edit", func() {
		nameEntry := widget.NewEntry()
		nameEntry.SetText(e.Name)
		dialog.ShowCustomConfirm("Edit Name", "Save", "Cancel",
			nameEntry,
			func(save bool) {
				if save && nameEntry.Text != "" {
					e.Name = nameEntry.Text
					nameLabel.SetText(e.Name)
					if onSave != nil {
						onSave()
					}
				}
			}, fyne.CurrentApp().Driver().AllWindows()[0])
	})

	nameContainerItems := []fyne.CanvasObject{nameLabel, editButton}
	if onExport != nil {
		exportBtn := widget.NewButton("Export to Bank", func() {
			onExport()
		})
		nameContainerItems = append(nameContainerItems, exportBtn)
	}

	if onImport != nil {
		importBtn := widget.NewButton("Import to Storage", func() {
			onImport()
		})
		nameContainerItems = append(nameContainerItems, importBtn)
	}

	if onRelease != nil {
		releaseBtn := widget.NewButton("Release", func() {
			onRelease()
		})
		nameContainerItems = append(nameContainerItems, releaseBtn)
	}

	nameContainer := container.NewHBox(nameContainerItems...)
	speciesInfo := e.Species + stellar
	infoLabel := widget.NewLabel(fmt.Sprintf(
		`%s | %s/%s | Lvl %d | HP %d/%d
Atk %d/%d | Def %d/%d | Spd %d
%s, %s, %s, %s | Emp: %s`,
		speciesInfo,
		getElementName(e.Element),
		getElementName(e.SubElement),
		e.CurrentLevel,
		e.Health,
		e.MaxHealth,
		e.PhysicalAttack,
		e.SpecialAttack,
		e.PhysicalDefense,
		e.SpecialDefense,
		e.Speed,
		e.Ability0Name,
		e.Ability1Name,
		e.Ability2Name,
		e.Ability3Name,
		e.EmpoweredAbilityName,
	))

	contentItems := []fyne.CanvasObject{nameContainer, infoLabel}

	content := container.NewVBox(contentItems...)

	return widget.NewCard("", "", content)
}

func createPlayerInfoCard(gameSave *GameSave, onSave func()) *widget.Card {
	nameLabel := widget.NewLabel(fmt.Sprintf("Name: %s", gameSave.ActivePlayerData.Name))
	elementLabel := widget.NewLabel(fmt.Sprintf("Spirit Element: %s", getElementName(gameSave.ActivePlayerData.SpiritElement)))

	genderOptions := []string{"Male", "Female"}
	selectedGender := "Male"
	if !gameSave.ActivePlayerData.IsMaleCharacter {
		selectedGender = "Female"
	}

	genderSelect := widget.NewSelect(genderOptions, func(selected string) {
		gameSave.ActivePlayerData.IsMaleCharacter = (selected == "Male")
		if onSave != nil {
			onSave()
		}
	})
	genderSelect.SetSelected(selectedGender)

	genderContainer := container.NewHBox(
		widget.NewLabel("Gender:"),
		genderSelect,
	)

	cardContent := container.NewVBox(
		nameLabel,
		elementLabel,
		genderContainer,
	)

	return widget.NewCard("Player Info", "", cardContent)
}

func createTeamTab(gameSave *GameSave, bank *Bank, onSave func(), onBankUpdate func(), myWindow fyne.Window) fyne.CanvasObject {
	elestrals := []*Elestral{
		gameSave.ActivePlayerData.Character0,
		gameSave.ActivePlayerData.Character1,
		gameSave.ActivePlayerData.Character2,
		gameSave.ActivePlayerData.Character3,
	}

	var cards []fyne.CanvasObject
	playerInfo := createPlayerInfoCard(gameSave, onSave)
	cards = append(cards, playerInfo)

	for _, e := range elestrals {
		elestral := e
		onExport := func() {
			if elestral != nil && elestral.Species != "" {
				elesCopy := *elestral
				bank.Elestrals = append(bank.Elestrals, &elesCopy)
				if onBankUpdate != nil {
					onBankUpdate()
				}
				dialog.ShowInformation("Export Successful",
					fmt.Sprintf("%s has been exported to the bank!", elestral.Name), myWindow)
			}
		}
		if card := createElestralCard(e, onSave, onExport, nil, nil); card != nil {
			cards = append(cards, card)
		}
	}

	content := container.NewVBox(cards...)
	return container.NewVScroll(content)
}

func createStorageTab(gameSave *GameSave, bank *Bank, onSave func(), onBankUpdate func(), myWindow fyne.Window) fyne.CanvasObject {
	var boxTabs []*container.TabItem
	for i, box := range gameSave.StorageBoxes {
		var cards []fyne.CanvasObject

		headerLabel := widget.NewLabel(fmt.Sprintf("Storage Box %d - %d Elestrals", i+1, len(box.Entries)))
		headerLabel.TextStyle = fyne.TextStyle{Bold: true}
		cards = append(cards, headerLabel)

		for _, entry := range box.Entries {
			elestral := entry.CharacterData
			onExport := func() {
				if elestral != nil && elestral.Species != "" {
					elesCopy := *elestral
					bank.Elestrals = append(bank.Elestrals, &elesCopy)
					if onBankUpdate != nil {
						onBankUpdate()
					}
					dialog.ShowInformation("Export Successful",
						fmt.Sprintf("%s has been exported to the bank!", elestral.Name), myWindow)
				}
			}
			if card := createElestralCard(entry.CharacterData, onSave, onExport, nil, nil); card != nil {
				cards = append(cards, card)
			}
		}

		content := container.NewVBox(cards...)
		scrollContainer := container.NewVScroll(content)
		boxTab := container.NewTabItem(fmt.Sprintf("Box %d", i+1), scrollContainer)
		boxTabs = append(boxTabs, boxTab)
	}

	return container.NewAppTabs(boxTabs...)
}

func findFirstAvailableSlot(gameSave *GameSave) (int, int, bool) {
	for boxIdx, box := range gameSave.StorageBoxes {
		for entryIdx, entry := range box.Entries {
			if entry.CharacterData == nil || entry.CharacterData.Species == "" {
				return boxIdx, entryIdx, true
			}
		}
	}
	return -1, -1, false
}

func createBankTab(gameSave *GameSave, bank *Bank, onSave func(), onBankUpdate func(), myWindow fyne.Window) fyne.CanvasObject {
	var cards []fyne.CanvasObject

	headerLabel := widget.NewLabel(fmt.Sprintf("Elestral Bank - %d Elestrals", len(bank.Elestrals)))
	headerLabel.TextStyle = fyne.TextStyle{Bold: true}
	cards = append(cards, headerLabel)

	for i, elestral := range bank.Elestrals {
		index := i
		eles := elestral

		onImport := func() {
			boxIdx, entryIdx, found := findFirstAvailableSlot(gameSave)
			if !found {
				dialog.ShowError(fmt.Errorf("no available slots in storage boxes"), myWindow)
				return
			}

			elesCopy := *eles
			gameSave.StorageBoxes[boxIdx].Entries[entryIdx].CharacterData = &elesCopy

			if onSave != nil {
				onSave()
			}
			if onBankUpdate != nil {
				onBankUpdate()
			}

			dialog.ShowInformation("Import Successful",
				fmt.Sprintf("%s has been imported to Storage Box %d!", eles.Name, boxIdx+1), myWindow)
		}

		onRelease := func() {
			dialog.ShowConfirm("Release Elestral",
				fmt.Sprintf("Are you sure you want to release %s? This cannot be undone.", eles.Name),
				func(confirm bool) {
					if !confirm {
						return
					}

					bank.Elestrals = append(bank.Elestrals[:index], bank.Elestrals[index+1:]...)
					if onBankUpdate != nil {
						onBankUpdate()
					}

					dialog.ShowInformation("Released",
						fmt.Sprintf("%s has been released.", eles.Name), myWindow)
				}, myWindow)
		}

		if card := createElestralCard(elestral, nil, nil, onImport, onRelease); card != nil {
			cards = append(cards, card)
		}
	}

	if len(bank.Elestrals) == 0 {
		emptyLabel := widget.NewLabel("No Elestrals in bank. Export Elestrals from your team or storage to add them here.")
		emptyLabel.Wrapping = fyne.TextWrapWord
		cards = append(cards, emptyLabel)
	}

	content := container.NewVBox(cards...)
	return container.NewVScroll(content)
}

func getSettingsFilePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, ".pbank.json"), nil
}

func loadSettings() (*Settings, error) {
	settingsPath, err := getSettingsFilePath()
	if err != nil {
		return &Settings{}, nil
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Settings{}, nil
		}
		return nil, err
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

func saveSettings(settings *Settings) error {
	settingsPath, err := getSettingsFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, data, 0644)
}

func getBankFilePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "pbank_store.json"), nil
}

func loadBank() (*Bank, error) {
	bankPath, err := getBankFilePath()
	if err != nil {
		return &Bank{Elestrals: []*Elestral{}}, nil
	}

	data, err := os.ReadFile(bankPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Bank{Elestrals: []*Elestral{}}, nil
		}
		return nil, err
	}

	var bank Bank
	if err := json.Unmarshal(data, &bank); err != nil {
		return nil, err
	}

	return &bank, nil
}

func saveBank(bank *Bank) error {
	bankPath, err := getBankFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(bank, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(bankPath, data, 0644)
}

func getDefaultSavePath(settings *Settings) string {
	if settings.CustomSavePath != "" {
		if _, err := os.Stat(settings.CustomSavePath); err == nil {
			return settings.CustomSavePath
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	var savePath string
	switch runtime.GOOS {
	case "windows":
		savePath = filepath.Join(homeDir, "AppData", "LocalLow", "DefaultCompany", "ElestralsAwakened-Playtest", "gamesave.json")
	case "darwin":
		savePath = filepath.Join(homeDir, "Library", "Application Support", "CrossOver", "Bottles", "Steam", "drive_c", "users", "crossover", "AppData", "LocalLow", "DefaultCompany", "ElestralsAwakened-Playtest", "gamesave.json")
	default:
		return ""
	}

	if _, err := os.Stat(savePath); err == nil {
		return savePath
	}

	return ""
}

func loadGameSave(filePath string) (*GameSave, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var gameSave GameSave
	if err := json.Unmarshal(data, &gameSave); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return &gameSave, nil
}

func saveGameSave(filePath string, gameSave *GameSave) error {
	if gameSave.ActivePlayerData.MaxSp == 0 {
		gameSave.ActivePlayerData.MaxSp = 100
		gameSave.ActivePlayerData.CurrentSp = 100
		gameSave.ActivePlayerData.BondMeter = 3
	}

	data, err := json.MarshalIndent(gameSave, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	return nil
}

func backupSaveFile(sourcePath string, myWindow fyne.Window) {
	now := time.Now()
	defaultName := fmt.Sprintf("awakened_backup_%d_%02d_%02d.json",
		now.Year(), now.Month(), now.Day())

	homeDir, _ := os.UserHomeDir()
	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()
		destPath := writer.URI().Path()
		if err := copyFile(sourcePath, destPath); err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		dialog.ShowInformation("Backup Successful",
			fmt.Sprintf("Save file backed up to:\n%s", destPath), myWindow)
	}, myWindow)

	saveDialog.SetFileName(defaultName)
	if homeURI, err := storage.ListerForURI(storage.NewFileURI(homeDir)); err == nil {
		saveDialog.SetLocation(homeURI)
	}
	saveDialog.Show()
}

func restoreSaveFile(destPath string, myWindow fyne.Window) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()

		sourcePath := reader.URI().Path()

		dialog.ShowConfirm("Confirm Restore",
			fmt.Sprintf("This will overwrite your current save file at:\n%s\n\nAre you sure?", destPath),
			func(confirm bool) {
				if !confirm {
					return
				}

				if err := copyFile(sourcePath, destPath); err != nil {
					dialog.ShowError(err, myWindow)
					return
				}

				dialog.ShowInformation("Restore Successful",
					"Save file has been restored successfully!", myWindow)
			}, myWindow)
	}, myWindow)
}

func changeDefaultSaveLocation(settings *Settings, myWindow fyne.Window) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()

		newPath := reader.URI().Path()
		settings.CustomSavePath = newPath
		if err := saveSettings(settings); err != nil {
			dialog.ShowError(fmt.Errorf("error saving settings: %w", err), myWindow)
			return
		}

		dialog.ShowInformation("Default Save Location Updated",
			"The default save location has been updated.\n\nPlease restart the application for changes to take effect.", myWindow)
	}, myWindow)
}

func displayGameSave(filePath string, myWindow fyne.Window, bank *Bank, tabs **container.AppTabs) {
	gameSave, err := loadGameSave(filePath)
	if err != nil {
		dialog.ShowError(err, myWindow)
		return
	}

	onSave := func() {
		err := saveGameSave(filePath, gameSave)
		if err != nil {
			dialog.ShowError(err, myWindow)
		}
	}

	var onBankUpdate func()
	onBankUpdate = func() {
		err := saveBank(bank)
		if err != nil {
			dialog.ShowError(fmt.Errorf("error saving bank: %w", err), myWindow)
		}

		if *tabs != nil {
			(*tabs).Items[0].Content = createTeamTab(gameSave, bank, onSave, onBankUpdate, myWindow)
			(*tabs).Items[1].Content = createStorageTab(gameSave, bank, onSave, onBankUpdate, myWindow)
			(*tabs).Items[2].Content = createBankTab(gameSave, bank, onSave, onBankUpdate, myWindow)
			(*tabs).Refresh()
		}
	}

	*tabs = container.NewAppTabs(
		container.NewTabItem("Team", createTeamTab(gameSave, bank, onSave, onBankUpdate, myWindow)),
		container.NewTabItem("Storage", createStorageTab(gameSave, bank, onSave, onBankUpdate, myWindow)),
		container.NewTabItem("Bank", createBankTab(gameSave, bank, onSave, onBankUpdate, myWindow)),
	)

	myWindow.SetContent(*tabs)
}

func showFilePickerDialog(myWindow fyne.Window, bank *Bank, tabs **container.AppTabs) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()

		filePath := reader.URI().Path()
		displayGameSave(filePath, myWindow, bank, tabs)
	}, myWindow)
}

func main() {
	myApp := app.NewWithID("com.primaryartemis.pandorasbank")
	myWindow := myApp.NewWindow("Pandora's Bank")

	settings, err := loadSettings()
	if err != nil {
		dialog.ShowError(fmt.Errorf("error loading settings: %w", err), myWindow)
		settings = &Settings{}
	}

	bank, err := loadBank()
	if err != nil {
		dialog.ShowError(fmt.Errorf("error loading bank: %w", err), myWindow)
		bank = &Bank{Elestrals: []*Elestral{}}
	}

	var tabs *container.AppTabs
	defaultSavePath := getDefaultSavePath(settings)
	var welcomeContent *fyne.Container
	if defaultSavePath != "" {
		welcomeLabel := widget.NewLabel("Welcome to Pandora's Bank!\n\nA game save file was found at the default location.")
		welcomeLabel.Alignment = fyne.TextAlignCenter

		openDefaultButton := widget.NewButton("Open Default Save", func() {
			displayGameSave(defaultSavePath, myWindow, bank, &tabs)
		})

		browseButton := widget.NewButton("Browse for Different Save...", func() {
			showFilePickerDialog(myWindow, bank, &tabs)
		})

		backupButton := widget.NewButton("Backup Save", func() {
			backupSaveFile(defaultSavePath, myWindow)
		})

		restoreButton := widget.NewButton("Restore Save", func() {
			restoreSaveFile(defaultSavePath, myWindow)
		})

		changeDefaultButton := widget.NewButton("Change Default Save Location", func() {
			changeDefaultSaveLocation(settings, myWindow)
		})

		welcomeContent = container.NewVBox(
			widget.NewLabel(""),
			welcomeLabel,
			widget.NewLabel(""),
			container.NewCenter(openDefaultButton),
			container.NewCenter(browseButton),
			widget.NewLabel(""),
			widget.NewSeparator(),
			widget.NewLabel(""),
			container.NewCenter(container.NewHBox(backupButton, restoreButton)),
			widget.NewLabel(""),
			container.NewCenter(changeDefaultButton),
		)
	} else {
		welcomeLabel := widget.NewLabel("Welcome to Pandora's Bank!\n\nPlease select a game save file to view your Elestrals.")
		welcomeLabel.Alignment = fyne.TextAlignCenter

		openButton := widget.NewButton("Open Game Save", func() {
			showFilePickerDialog(myWindow, bank, &tabs)
		})

		changeDefaultButton := widget.NewButton("Set Default Save Location", func() {
			changeDefaultSaveLocation(settings, myWindow)
		})

		welcomeContent = container.NewVBox(
			widget.NewLabel(""),
			welcomeLabel,
			widget.NewLabel(""),
			container.NewCenter(openButton),
			widget.NewLabel(""),
			widget.NewSeparator(),
			widget.NewLabel(""),
			container.NewCenter(changeDefaultButton),
		)
	}

	myWindow.SetContent(welcomeContent)
	myWindow.Resize(fyne.NewSize(600, 800))
	myWindow.ShowAndRun()
}
