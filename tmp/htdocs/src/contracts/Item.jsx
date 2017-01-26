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
    return (
      <div class="item-cell">
        <p class="lbl"><span class="index">{this.props.index}</span> {release.tender.id}</p>
        <p><a onClick={this.handleClick}>{release.tender.description}</a></p>
        <div>
          <div class="left">
            <p class="lbl">Monto</p>
            <p><strong>{formatAmount( release.tender.value.amount || 0 )}</strong></p>
            <p class="lbl">Proveedor</p>
            <p>{release.awards[0].suppliers[0].name}</p>
          </div>
          <div class="right">
            <p class="lbl">Comprador</p>
            <p>{release.buyer.name}</p>
            <span class="btn-black active" onClick={this.handleClick}>Ver contrato</span>
          </div>
          <div class="clear"></div>
        </div>
      </div>
    );
  }
}

export default Item;
