export const REDEEM_UPDATE_GIFT_CARD_FIELD_1 =
  'REDEEM_UPDATE_GIFT_CARD_FIELD_1';
export const REDEEM_UPDATE_GIFT_CARD_FIELD_2 =
  'REDEEM_UPDATE_GIFT_CARD_FIELD_2';
export const REDEEM_UPDATE_GIFT_CARD_FIELD_3 =
  'REDEEM_UPDATE_GIFT_CARD_FIELD_3';

export const REDEEM_ADD_GIFT_CARD = 'REDEEM_ADD_GIFT_CARD';
export const REDEEM_REMOVE_GIFT_CARD = 'REDEEM_REMOVE_GIFT_CARD';

export const REDEEM_UPDATE_GIFT_CARD_ENTRY_ERROR =
  'REDEEM_UPDATE_GIFT_CARD_ENTRY_ERROR';

export const REDEEM_REDEEM_GIFT_CARDS = 'REDEEM_GIFT_CARDS';
export const REDEEM_REDEEM_GIFT_CARDS_SUCCESS =
  'REDEEM_REDEEM_GIFT_CARDS_SUCCESS';

export function updateGiftCardField1(content) {
  return {
    type: REDEEM_UPDATE_GIFT_CARD_FIELD_1,
    content
  };
}

export function updateGiftCardField2(content) {
  return {
    type: REDEEM_UPDATE_GIFT_CARD_FIELD_2,
    content
  };
}

export function updateGiftCardField3(content) {
  return {
    type: REDEEM_UPDATE_GIFT_CARD_FIELD_3,
    content
  };
}

export function addGiftCard(claimCode) {
  return {
    type: REDEEM_ADD_GIFT_CARD,
    claimCode
  };
}

export function removeGiftCard(claimCode) {
  return {
    type: REDEEM_REMOVE_GIFT_CARD,
    claimCode
  };
}

export function updateGiftCardEntryError(errorText) {
  return {
    type: REDEEM_UPDATE_GIFT_CARD_ENTRY_ERROR,
    errorText
  };
}

export function redeemGiftCards(id, giftCards) {
  return {
    type: REDEEM_REDEEM_GIFT_CARDS,
    giftCards,
    id
  };
}

export function redeemGiftCardsSuccess(resp) {
  return {
    type: REDEEM_REDEEM_GIFT_CARDS_SUCCESS,
    resp
  };
}
