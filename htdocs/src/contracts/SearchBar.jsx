import React from 'react';
import { formatDate } from '../helpers.js';

class SearchBar extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      filter: null,
      value: null,
      limit: 30
    }
  }

  componentDidMount() {
    this.ui = {}
    this.ui.form = $('form#queryForm');
    this.ui.bullets = this.ui.form.find('div.bullets span');
    this.ui.input = this.ui.form.find( 'input#query' );

    // Setup calendar
    this.ui.form.find("span[data-filter='date']").datepicker({
      clearBtn: true,
      assumeNearbyYear: true,
      format: 'mm/dd/yyyy',
      language: 'es',
      maxViewMode: 2,
      multidate: 2,
      todayHighlight: true
    }).on({
      'hide': (e) => {
        e.dates.sort(function(a, b) {
          return new Date(a).getTime() - new Date(b).getTime();
        });
        let lbl = formatDate(e.dates[0], 'MMMM Do YYYY');
        let val = formatDate(e.dates[0], 'MM-DD-YYYY');
        if( e.dates.length > 1 ) {
          lbl += ' a ' + formatDate(e.dates[1], 'MMMM Do YYYY')
          val += '|' + formatDate(e.dates[1], 'MM-DD-YYYY');
        }

        // Update state
        this.setState({
          value: val
        });

        // Update UI
        this.ui.input.val( lbl );
        this.ui.input.focus();
      },
      'changeDate': () => {
        // Visually mark the selected range
        if( $('div.datepicker-days tbody tr td.active').length == 2 ) {
          let mark = false;
          let days = $('div.datepicker-days tbody td.active');
          $('div.datepicker-days tbody td').each(function(i, d) {
            if( d == days[0] ) {
              mark = true;
            }
            if( mark ) {
              $(d).addClass('range');
            }
            if( d == days[1] ) {
              mark = false;
            }
          });
        }
      }
    });

    // Setup slider
    this.ui.form.find("span[data-filter='amount']").popover({
      html: true,
      title: 'Seleccione el rango a utilizar como filtro (MXN)',
      content: '<b>$0</b><input id="amountSlider" type="text" /><b>$100,000,000</b>',
      placement: 'bottom',
      trigger: 'focus'
    }).on( 'shown.bs.popover', () => {
      $("#amountSlider").slider({
        step: 50000,
        min: 0,
        max: 100000000,
        value: [20000000,80000000],
        formatter: function( value ) {
          if( Array.isArray( value ) ) {
            var lbl = '$' + value[0].toLocaleString() + ' a ' + '$' + value[1].toLocaleString();
            return lbl;
          }
          return '';
        }
      }).on( 'slide', (e) => {
        this.setState({
          value: e.value.join('|')
        });
        this.ui.input.val( '$' + e.value[0].toLocaleString() + ' a ' + '$' + e.value[1].toLocaleString() );
      });
    });

    // Update state when filter value changes
    this.ui.input.on( 'keyup', () => {
      this.setState({
        value: this.ui.input.val()
      })
    });

    // Handle filter selection
    this.ui.bullets.on( 'click', (function( e ) {
      // Toggle active
      let target = $( e.target );
      this.ui.bullets.removeClass( 'active' );
      target.addClass( 'active' );

      // Update state
      this.setState({
        filter: target.data('filter')
      });

      switch (target.data( 'filter' )) {
        case 'date':
          target.datepicker( 'show' );
          break;
        case 'amount':
          target.popover( 'toggle' );
          break;
      }
    }).bind(this) );

    // Handle form submit
    this.ui.form.on( 'submit', (function( e ) {
      e.preventDefault();
      this.props.onSubmit(this.state);
    }).bind(this));

    // Clear filters
    this.ui.form.on( 'reset', (e) => {
      e.preventDefault();
      this.setState({
        filter: null,
        value: null,
        limit: 30
      });
      this.ui.input.val( '' );
    });
  }

  render() {
    return (
      <div className="inner-row">
        <div className="row">
          <div className="col-md-12">
            <h2 className="mobile-search">Buscador de Contratos</h2>
            <p className="mobile-search">Encuentra contratos o procedimientos de contratación registrados en Testigo Social 2.0 haciendo uso de los distintos filtros de búsqueda disponibles.</p>
            <form id="queryForm">
              <div className="input-group hide-contracts">
                <input type="text" className="form-control" id="query" name="query" placeholder="Buscar..." />
                <span className="input-group-btn">
                  <button className="btn" type="reset">Limpiar</button>
                  <button className="btn btn-primary" type="submit">Buscar</button>
                </span>
              </div>
              <div className="bullets">
                <span className="btn-black active" data-filter="date">Fecha</span>
                <span className="btn-black" data-filter="amount">Monto</span>
                <span className="btn-black" data-filter="buyer">Comprador</span>
                <span className="btn-black" data-filter="provider">Proveedor</span>
                <span className="btn-black" data-filter="procedureType">Tipo de Procedimiento</span>
                <span className="btn-black" data-filter="procedureNumber">No. de Procedimiento</span>
                <span className="btn-black" data-filter="contractNumber">No. de Contrato</span>
              </div>

              <div className="input-group mobile-contracts">
                <input type="text" className="form-control" id="query" name="query" placeholder="Buscar" />
                <button className="btn btn-primary btn-mobile" type="submit">Buscar</button>
              </div>
            </form>
          </div>
        </div>
      </div>
    );
  }
}

export default SearchBar;
