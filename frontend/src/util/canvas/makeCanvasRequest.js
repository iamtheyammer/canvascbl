import axios from 'axios';
import getUrlPrefix from '../getUrlPrefix';

export default (
  path,
  token,
  subdomain = 'canvas',
  query = {},
  method = 'get'
) =>
  axios({
    method,
    url: `${getUrlPrefix}/api/canvas/${path}`,
    headers: {
      'X-Canvas-Token': token,
      'X-Canvas-Subdomain': subdomain
    },
    params: query
  });
