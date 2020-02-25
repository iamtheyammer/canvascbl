import {
  SETTINGS_GET_NOTIFICATION_SETTINGS_AND_TYPES_ERROR,
  SETTINGS_GOT_NOTIFICATION_SETTINGS_AND_TYPES,
  SETTINGS_TOGGLE_NOTIFICATION_TYPE_ERROR
} from '../actions/settings';

function settings(state = {}, action) {
  switch (action.type) {
    case SETTINGS_GET_NOTIFICATION_SETTINGS_AND_TYPES_ERROR:
      return {
        ...state,
        getNotificationSettingsAndTypesError: action.e
      };
    case SETTINGS_GOT_NOTIFICATION_SETTINGS_AND_TYPES:
      return {
        ...state,
        notifications: {
          ...state.notifications,
          settings: action.settings,
          types: action.types
        }
      };
    case SETTINGS_TOGGLE_NOTIFICATION_TYPE_ERROR:
      return {
        ...state,
        notifications: {
          ...state.notifications,
          toggles: {
            ...state.notifications.toggles,
            [action.typeId]: action.e
          }
        }
      };
    default:
      return state;
  }
}

export default settings;
