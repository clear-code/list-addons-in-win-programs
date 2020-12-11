/*
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
'use strict';

const mAddons = new Map();

browser.management.getAll().then(async addons => {
  addons.forEach(addToRegistry);
  const ids = await getRegisteredAddonIds();
  for (const id of ids) {
    if (!mAddons.has(id))
      removeFromRegistry(id);
  }
});

browser.management.onInstalled.addListener(addon => {
  console.log('onInstalled: ', addon);
  addToRegistry(addon);
});

// This won't be called when the addon is uninstalled via about:addons. Why?
browser.management.onUninstalled.addListener(addon => {
  console.log('onUninstalled: ', addon);
  removeFromRegistry(addon.id);
});

// This never been called...
browser.management.onDisabled.addListener(addon => {
  console.log('onDisabled: ', addon);
  if (addon.id != browser.runtime.id)
    return;

  for (const id of mAddons.keys()) {
    removeFromRegistry(id);
  }
});

function addToRegistry(addon) {
  console.log('addToRegistry: ', addon);
  mAddons.set(addon.id, addon);
}

function removeFromRegistry(id) {
  console.log('removeFromRegistry: ', id);
  mAddons.delete(id);
}

async function getRegisteredAddonIds() {
  return [];
}
