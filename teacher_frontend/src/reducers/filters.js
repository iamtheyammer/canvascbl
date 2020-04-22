import {
  FILTERS_CLEAR,
  FILTERS_UPDATE_NAME,
  FILTERS_UPDATE_NAME_TYPE
} from '../actions/filters';

export default function filters(state = {}, action) {
  switch (action.type) {
    case FILTERS_CLEAR:
      return {};
    case FILTERS_UPDATE_NAME:
      return {
        ...state,
        name: action.newName
      };
    case FILTERS_UPDATE_NAME_TYPE:
      return {
        ...state,
        nameType: action.newType
      };
    default:
      return state;
  }
}
