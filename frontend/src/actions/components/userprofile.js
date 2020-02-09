export const USERPROFILE_SET_GET_AUTHORIZED_APPS_ID =
  'USERPROFILE_SET_GET_AUTHORIZED_APPS_ID';

export function setGetAuthorizedAppsId(id) {
  return {
    type: USERPROFILE_SET_GET_AUTHORIZED_APPS_ID,
    id
  };
}
