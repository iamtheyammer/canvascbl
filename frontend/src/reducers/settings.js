import {
  SETTINGS_GET_NOTIFICATION_SETTINGS_AND_TYPES_ERROR,
  SETTINGS_GOT_NOTIFICATION_SETTINGS_AND_TYPES,
  SETTINGS_TOGGLE_NOTIFICATION_TYPE_ERROR,
  SETTINGS_TOGGLED_SHOW_HIDDEN_COURSES
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
          toggleErrors: {
            ...state.notifications.toggleErrors,
            [action.typeId]: action.e
          }
        }
      };
    case SETTINGS_TOGGLED_SHOW_HIDDEN_COURSES:
      return {
        ...state,
        showHiddenCourses: action.toggle
      };
    default:
      return state;
  }
}

export default settings;
