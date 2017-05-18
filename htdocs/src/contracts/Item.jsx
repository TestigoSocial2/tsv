import React from 'react';
import { formatAmount } from '../helpers.js';

class Item extends React.Component {
  constructor(props) {
    super(props);
    this.handleClick = this.handleClick.bind(this);
  }

  handleClick(e) {
    e.preventDefault();
    this.props.onClick( this.props.contract );
  }

  render() {
    let release = this.props.contract.releases[0];
    let supplier = '-';
    if( release.awards[0].supplier ) {
      supplier = release.awards[0].supplier.name;
    }

    return (
      <div className="item-cell">
        <p className="lbl"><span className="index">{this.props.index}</span> {release.tender.id}</p>
        <p><a onClick={this.handleClick}>{release.tender.description}</a></p>
        <div>
          <div className="left">
            <p className="lbl">Monto</p>
            <p><strong>{formatAmount( release.planning.budget.amount.amount || 0 )}</strong></p>
            <p className="lbl">Proveedor</p>
            <p>{supplier}</p>
          </div>
          <div className="right">
            <p className="lbl">Comprador</p>
            <p>{release.buyer.name}</p>
            <span className="btn-black active" onClick={this.handleClick}>Ver contrato</span>
          </div>
          <div className="clear"></div>
        </div>
      </div>
    );
  }
}

export default Item;
