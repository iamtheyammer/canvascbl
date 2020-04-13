import { GRADES_SWITCH_VIEW_TYPE } from '../../actions/components/grades';

export default function grades(state = {}, action) {
  switch (action.type) {
    case GRADES_SWITCH_VIEW_TYPE:
      return {
        ...state,
        viewType: action.newTypeName
      };
    default:
      return state;
  }
}
