function Controller() {
	installer.autoRejectMessageBoxes();
}

Controller.prototype.WelcomePageCallback = function() {
	// click delay here because the next button is initially disabled for ~1 second
	gui.clickButton(buttons.NextButton, 3000);
}

Controller.prototype.CredentialsPageCallback = function() {
	gui.clickButton(buttons.NextButton);
}

Controller.prototype.IntroductionPageCallback = function() {
	gui.clickButton(buttons.NextButton);
}

Controller.prototype.TargetDirectoryPageCallback = function() {
	gui.currentPageWidget().TargetDirectoryLineEdit.setText("/home/user/Qt5.13.0");
	gui.clickButton(buttons.NextButton);
}

Controller.prototype.ComponentSelectionPageCallback = function() {
	var qtVersion = 'qt5.5130',
		widget = gui.currentPageWidget();
	
	widget.deselectAll();
	widget.selectComponent('qt.'+qtVersion+'.gcc_64');
	widget.selectComponent('qt.'+qtVersion+'.qtcharts');
	widget.selectComponent('qt.'+qtVersion+'.qtdatavis3d');
	gui.clickButton(buttons.NextButton);
}

Controller.prototype.LicenseAgreementPageCallback = function() {
	gui.currentPageWidget().AcceptLicenseRadioButton.setChecked(true);
	gui.clickButton(buttons.NextButton);
}

Controller.prototype.StartMenuDirectoryPageCallback = function() {
	gui.clickButton(buttons.NextButton);
}

Controller.prototype.ReadyForInstallationPageCallback = function() {
	gui.clickButton(buttons.NextButton);
}

Controller.prototype.PerformInstallationPageCallback = function() {
	gui.currentPageWidget().setAutomatedPageSwitchEnabled(true);
}

Controller.prototype.FinishedPageCallback = function() {
	var checkBoxForm = gui.currentPageWidget().LaunchQtCreatorCheckBoxForm;
	if (checkBoxForm && checkBoxForm.launchQtCreatorCheckBox) {
		checkBoxForm.launchQtCreatorCheckBox.checked = false;
	}

	gui.clickButton(buttons.FinishButton);
}