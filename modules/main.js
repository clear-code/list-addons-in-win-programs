/*
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

var registry = require('registry').registry;
var { AddonManager } = Cu.import('resource://gre/modules/AddonManager.jsm', {});
var { Services } = Cu.import('resource://gre/modules/Services.jsm', {});
var basePath = 'HKEY_CURRENT_USER\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall';

var installed = registry.getChildren(basePath);

console.log(installed.join('\n'));


AddonManager.getAllAddons(function(addons) {
  addons.forEach(function(addon) {
    if (addon.type !== 'extension')
      return;
    console.log(Services.appinfo.ID);
    console.log(addon.id);
    console.log(addon.name);
    console.log(addon.version);
    console.log(addon.type);
  });
});
