/*
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

var registry = require('registry').registry;
var { AddonManager } = Cu.import('resource://gre/modules/AddonManager.jsm', {});
var { Services } = Cu.import('resource://gre/modules/Services.jsm', {});
var { FileUtils } = Cu.import('resource://gre/modules/FileUtils.jsm', {});

var exePath = FileUtils.getFile("XREExeF", []).path;
var basePath = 'HKEY_CURRENT_USER\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall';

function log(aMessage) {
  console.log('[list-addons-in-win-programs] ' + aMessage);
}

log('Services.appinfo.ID: ' + Services.appinfo.ID);
log('Services.appinfo.name: ' + Services.appinfo.name);

function createRegistryKey(aAddon) {
  log('createRegistryKey');
  var key = basePath + '\\' + Services.appinfo.ID + '.' + aAddon.id;
  log('key: ' + key);
  return key;
}

function registerUninstallInfo(aKey, aAddon) {
  log('registerUninstallInfo');
  log('aKey: ' + aKey);
  log('aAddon: ' + aAddon.id);
  registry.setValue(aKey + '\\' + 'DisplayName', Services.appinfo.name + ': ' + aAddon.name);
  registry.setValue(aKey + '\\' + 'DisplayVersion', aAddon.version);
  registry.setValue(aKey + '\\' + 'UninstallString', exePath);
  registry.setValue(aKey + '\\' + 'DisplayIcon', exePath + ',0');
  registry.setValue(aKey + '\\' + 'Publisher', aAddon.creator.name);
}

AddonManager.getAllAddons(function(aAddons) {
  var currentInstalledAddonKeys = [];
  aAddons.forEach(function(aAddon) {
    if (aAddon.type !== 'extension')
      return;
    var key = createRegistryKey(aAddon);
    registerUninstallInfo(key, aAddon);
    currentInstalledAddonKeys.push(key);
  });
  log('currentInstalledAddonKeys: ' + JSON.stringify(currentInstalledAddonKeys));

  var installed = registry.getChildren(basePath);
  log('installed: ' + installed.join('\n'));
  installed.forEach(function(key) {
    log('installed key: ' + key);
    log('indexOf: ' + currentInstalledAddonKeys.indexOf(key));
    if (currentInstalledAddonKeys.indexOf(key) === -1) {
      log('registry.clear: ' + key);
      registry.clear(key);
    }
  });
});
