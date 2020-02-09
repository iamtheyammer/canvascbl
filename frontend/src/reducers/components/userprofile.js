import { USERPROFILE_SET_GET_AUTHORIZED_APPS_ID } from '../../actions/components/userprofile';

export default function userprofile(state = {}, action) {
  switch (action.type) {
    case USERPROFILE_SET_GET_AUTHORIZED_APPS_ID:
      return {
        ...state,
        getAuthorizedAppsId: action.id
      };
    default:
      return state;
  }
}
