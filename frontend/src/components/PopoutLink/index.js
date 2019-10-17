import React from 'react';
import * as PropTypes from 'prop-types';

function PopoutLink(props) {
  const { url, children } = props;
  return (
    <a target="_blank" rel="noopener noreferrer" href={url}>
      {children}
    </a>
  );
}

PopoutLink.propTypes = {
  url: PropTypes.string.isRequired
};

export default PopoutLink;
