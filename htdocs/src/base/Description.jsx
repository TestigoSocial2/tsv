import React from 'react';

// Section description block
class Description extends React.Component {
  render() {
    return(
      <div className={'row info bg-' + this.props.color}>
        <div className="inner-row">
          <h2>{this.props.title}</h2>
          <p>{this.props.content}</p>
        </div>
      </div>
    );
  }
}

export default Description;
