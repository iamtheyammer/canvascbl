import makePlusRequest from '../util/plus/makePlusRequest';
import { plusError } from './error';
import { startLoading, endLoading } from './loading';

export const PLUS_GOT_SESSION_INFORMATION = 'PLUS_GOT_SESSION_INFORMATION';

function gotSessionInformation(sessionInformation) {
  return {
    type: PLUS_GOT_SESSION_INFORMATION,
    sessionInformation
  };
}

export function getSessionInformation(id) {
  return async dispatch => {
    startLoading(id);
    try {
      const sessionRequest = await makePlusRequest('session');
      dispatch(gotSessionInformation(sessionRequest.data));
    } catch (e) {
      dispatch(plusError(id, e.res));
    }
    endLoading(id);
  };
}
