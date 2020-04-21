import { LOADING_UPDATE_NUMBER_OF_DOTS } from "../../actions/components/loading";

export default function loading(state = {}, action) {
  switch (action.type) {
    case LOADING_UPDATE_NUMBER_OF_DOTS:
      return {
        ...state,
        dots: action.numDots
      };
    default:
      return state;
  }
}
