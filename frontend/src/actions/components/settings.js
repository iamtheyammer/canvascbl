export const CSETTINGS_SET_GET_NOTIFICATION_SETTINGS_AND_TYPES_ID =
  'CSETTINGS_SET_GET_NOTIFICATION_SETTINGS_AND_TYPES_ID';

export const CSETTINGS_SET_TOGGLE_NOTIFICATION_STATUS_ID =
  'CSETTINGS_SET_TOGGLE_NOTIFICATION_STATUS_ID';

export function setGetNotificationSettingsAndTypesId(id) {
  return {
    type: CSETTINGS_SET_GET_NOTIFICATION_SETTINGS_AND_TYPES_ID,
    id
  };
}

export function setToggleNotificationStatusId(id, typeId) {
  return {
    type: CSETTINGS_SET_TOGGLE_NOTIFICATION_STATUS_ID,
    id,
    typeId
  };
}
