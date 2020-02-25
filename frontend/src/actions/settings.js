export const SETTINGS_GET_NOTIFICATION_SETTINGS_AND_TYPES =
  'SETTINGS_GET_NOTIFICATION_SETTINGS_AND_TYPES';
export const SETTINGS_GET_NOTIFICATION_SETTINGS_AND_TYPES_ERROR =
  'SETTINGS_GET_NOTIFICATION_SETTINGS_AND_TYPES_ERROR';
export const SETTINGS_GOT_NOTIFICATION_SETTINGS_AND_TYPES =
  'SETTINGS_GOT_NOTIFICATION_SETTINGS_AND_TYPES';

export const SETTINGS_TOGGLE_NOTIFICATION_TYPE =
  'SETTINGS_TOGGLE_NOTIFICATION_TYPE';
export const SETTINGS_TOGGLE_NOTIFICATION_TYPE_ERROR =
  'SETTINGS_TOGGLE_NOTIFICATION_TYPE_ERROR';
export const SETTINGS_TOGGLED_NOTIFICATION_TYPE =
  'SETTINGS_TOGGLED_NOTIFICATION_TYPE';

export function getNotificationSettingsAndTypes(id) {
  return {
    type: SETTINGS_GET_NOTIFICATION_SETTINGS_AND_TYPES,
    id
  };
}

export function getNotificationSettingsAndTypesError(e) {
  return {
    type: SETTINGS_GET_NOTIFICATION_SETTINGS_AND_TYPES_ERROR,
    e
  };
}

export function gotNotificationSettingsAndTypes(settings, types) {
  return {
    type: SETTINGS_GOT_NOTIFICATION_SETTINGS_AND_TYPES,
    settings,
    types
  };
}

export function toggleNotificationType(id, typeId, toggle) {
  return {
    type: SETTINGS_TOGGLE_NOTIFICATION_TYPE,
    id,
    typeId,
    toggle
  };
}

export function toggleNotificationTypeError(e, typeId) {
  return {
    type: SETTINGS_TOGGLE_NOTIFICATION_TYPE_ERROR,
    e,
    typeId
  };
}

export function toggledNotificationType(typeId) {
  return {
    type: SETTINGS_TOGGLED_NOTIFICATION_TYPE,
    typeId
  };
}
