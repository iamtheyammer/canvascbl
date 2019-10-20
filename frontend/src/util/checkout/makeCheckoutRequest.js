import axios from 'axios';
import getUrlPrefix from '../getUrlPrefix';

export default (
  path,
  query = {},
  method = 'get'
) =>
  axios({
    method,
    url: `${getUrlPrefix}/api/checkout/${path}`,
    params: query
  });
