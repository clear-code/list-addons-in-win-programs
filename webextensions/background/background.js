/*
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
'use strict';

const HOST_ID = 'com.clear_code.list_addons_in_win_programs_we_host';

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

async function addToRegistry(addon) {
  console.log('addToRegistry: ', addon);
  mAddons.set(addon.id, addon);
  try {
    const response = await sendToHost({
      type: 'register-addon',
      addon,
    });
    console.log('addToRegistry response: ', addon.id, response);
  }
  catch(error) {
    console.error(error);
  }
}

async function removeFromRegistry(id) {
  console.log('removeFromRegistry: ', id);
  mAddons.delete(id);
  try {
    const response = await sendToHost({
      type: 'unregister-addon',
      id,
    });
    console.log('removeFromRegistry response: ', id, response);
  }
  catch(error) {
    console.error(error);
  }
}

async function getRegisteredAddonIds() {
  /*
  try {
    const response = await sendToHost({
      type: 'list-registered-addons',
    });
  }
  catch(error) {
    console.error(error);
    return [];
  }
  */
  return [];
}

async function sendToHost(message) {
  try {
    const response = await browser.runtime.sendNativeMessage(HOST_ID, message);
    if (!response || typeof response != 'object')
      throw new Error(`invalid response: ${String(response)}`);
    return response;
  }
  catch(error) {
    console.log('Error: failed to get response for message', message, error);
    return null;
  }
}
