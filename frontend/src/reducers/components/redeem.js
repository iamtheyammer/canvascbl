import {
  REDEEM_ADD_GIFT_CARD,
  REDEEM_REDEEM_GIFT_CARDS,
  REDEEM_REDEEM_GIFT_CARDS_SUCCESS,
  REDEEM_REMOVE_GIFT_CARD,
  REDEEM_UPDATE_GIFT_CARD_ENTRY_ERROR,
  REDEEM_UPDATE_GIFT_CARD_FIELD_1,
  REDEEM_UPDATE_GIFT_CARD_FIELD_2,
  REDEEM_UPDATE_GIFT_CARD_FIELD_3
} from '../../actions/components/redeem';

export default function redeem(state = {}, action) {
  switch (action.type) {
    case REDEEM_UPDATE_GIFT_CARD_FIELD_1:
      return {
        ...state,
        giftCardField1: action.content
      };
    case REDEEM_UPDATE_GIFT_CARD_FIELD_2:
      return {
        ...state,
        giftCardField2: action.content
      };
    case REDEEM_UPDATE_GIFT_CARD_FIELD_3:
      return {
        ...state,
        giftCardField3: action.content
      };
    case REDEEM_ADD_GIFT_CARD:
      return {
        ...state,
        giftCards: state.giftCards
          ? [...state.giftCards, action.claimCode]
          : [action.claimCode],
        giftCardField1: '',
        giftCardField2: '',
        giftCardField3: '',
        giftCardEntryError: ''
      };
    case REDEEM_REMOVE_GIFT_CARD:
      return {
        ...state,
        giftCards:
          state.giftCards &&
          state.giftCards.filter((gc) => gc !== action.claimCode)
      };
    case REDEEM_UPDATE_GIFT_CARD_ENTRY_ERROR:
      return {
        ...state,
        giftCardEntryError: action.errorText
      };
    case REDEEM_REDEEM_GIFT_CARDS:
      return {
        ...state,
        redeemGiftCardsId: action.id
      };
    case REDEEM_REDEEM_GIFT_CARDS_SUCCESS:
      return {
        ...state,
        redeemGiftCardsResponse: action.resp
      };
    default:
      return state;
  }
}
