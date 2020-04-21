export const LOADING_UPDATE_NUMBER_OF_DOTS = 'LOADING_UPDATE_NUMBER_OF_DOTS';

export function updateNumberOfDots(numDots) {
  return {
    type: LOADING_UPDATE_NUMBER_OF_DOTS,
    numDots,
  };
}
