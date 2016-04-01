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

function log(message) {
  console.log('[list-addons-in-win-programs] ' + message);
}

var installed = registry.getChildren(basePath);

log('installed: ' + installed.join('\n'));

log('Services.appinfo.ID: ' + Services.appinfo.ID);
log('Services.appinfo.name: ' + Services.appinfo.name);

function writeUninstallInfo(aAddon) {
  log('writeUninstallInfo');
  log('aAddon: ' + aAddon.id);
  var key = basePath + '\\' + Services.appinfo.ID + '.' + aAddon.id + '\\';
  log('key: ' + key);
  registry.setValue(key + 'DisplayName', Services.appinfo.name + ': ' + aAddon.name);
  registry.setValue(key + 'DisplayVersion', aAddon.version);
  registry.setValue(key + 'UninstallString', exePath);
  registry.setValue(key + 'DisplayIcon', exePath + ',0');
  registry.setValue(key + 'Publisher', aAddon.creator.name);
}

AddonManager.getAllAddons(function(addons) {
  addons.forEach(function(addon) {
    if (addon.type !== 'extension')
      return;
    writeUninstallInfo(addon);
  });
});
