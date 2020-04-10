import { COURSESETTINGS_SET_TOGGLE_COURSE_VISIBILITY_ID } from '../../actions/components/coursesettings';

export default function coursesettings(state = {}, action) {
  switch (action.type) {
    case COURSESETTINGS_SET_TOGGLE_COURSE_VISIBILITY_ID:
      return {
        ...state,
        toggleVisibilityIds: {
          ...state.toggleVisibilityIds,
          [action.courseId]: action.id
        }
      };
    default:
      return state;
  }
}
