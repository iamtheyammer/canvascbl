import React, { useEffect } from 'react';
import * as PropTypes from 'prop-types';
import { v4 } from 'uuid';
import { Icon } from 'antd';
import { ReactComponent as PopOutIcon } from '../../assets/pop_out.svg';
import { trackExternalLinkClick } from '../../util/tracking';

function PopoutLink(props) {
  const { url, children, addIcon, id, tracking } = props;

  const elId = id ? id : v4();

  useEffect(() => {
    // this ensures that the element exists before sending officially
    // tracking it
    try {
      if (tracking) {
        trackExternalLinkClick(
          elId,
          url,
          tracking.destinationName,
          tracking.destinationType,
          tracking.via
        );
      }
    } catch (e) {}
    // we want this to run only once.
    // eslint-disable-next-line
  }, []);

  return (
    <a target="_blank" rel="noopener noreferrer" href={url} id={elId}>
      {children}
      {addIcon && ' '}
      {addIcon && <Icon component={PopOutIcon} />}
    </a>
  );
}

PopoutLink.propTypes = {
  // url to link to
  url: PropTypes.string.isRequired,
  // whether to add a space and the PopOutIcon after your content (children)
  addIcon: PropTypes.bool,
  // anything uniquely identifying this link
  id: PropTypes.any,
  // info to set up a tracking link
  tracking: PropTypes.shape({
    destinationName: PropTypes.string.isRequired,
    destinationType: PropTypes.string,
    via: PropTypes.string.isRequired
  })
};

export default PopoutLink;
