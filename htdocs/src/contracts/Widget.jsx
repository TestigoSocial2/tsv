import React from 'react';
import Item from './Item.jsx';

class Widget extends React.Component {
  constructor(props) {
    super(props);
    this.itemSelected = this.itemSelected.bind(this);
    this.state = {
      items: []
    }
  }

  componentDidMount() {
    $.ajax({
      type: "POST",
      url: "/query",
      data: this.props.query,
      success: (res) => {
        this.setState({items: res});
      }
    });
  }

  itemSelected(contract) {
    this.props.onSelection(contract);
  }

  render() {
    return(
      <div className="col-md-6">
        <h2 className="block-title">{this.props.title}</h2>
        <p>{this.props.description}</p>
        {
          this.state.items.map((contract, i) => (
            <Item
              index={i + 1}
              onClick={this.itemSelected}
              key={contract.releases[0].id}
              contract={contract} />
          ))
        }
      </div>
    );
  }
}

export default Widget;