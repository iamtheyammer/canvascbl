import {
  CSETTINGS_SET_GET_NOTIFICATION_SETTINGS_AND_TYPES_ID,
  CSETTINGS_SET_TOGGLE_NOTIFICATION_STATUS_ID
} from '../../actions/components/settings';

function settings(state = {}, action) {
  switch (action.type) {
    case CSETTINGS_SET_GET_NOTIFICATION_SETTINGS_AND_TYPES_ID:
      return {
        ...state,
        getNotificationSettingsAndTypesId: action.id
      };
    case CSETTINGS_SET_TOGGLE_NOTIFICATION_STATUS_ID:
      return {
        ...state,
        toggleIds: {
          ...state.toggleIds,
          [action.typeId]: action.id
        }
      };
    default:
      return state;
  }
}

export default settings;
