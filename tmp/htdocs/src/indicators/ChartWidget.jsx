import React from 'react';
import Chart from 'chart.js'

class ChartWidget extends React.Component {
  constructor(props) {
    super(props);
    this.chart = null;
  }

  render() {
    if( this.props.data.hasOwnProperty('limited') ) {
      // Build chart
      this.chart = new Chart( document.getElementById(this.props.id), {
        type: "pie",
        data: this.props.reducer(this.props.data),
        options: {
          responsive: true,
          responsiveAnimationDuration: 500,
          padding: 10
        }
      });
    }

    return (
      <div className="chart-widget">
        <h2 className="block-title">{this.props.title}</h2>
        <div className="bg-gray">
          <p>{this.props.description}</p>
        </div>
        <div className="chart">
          <canvas id={this.props.id}
            className="dataChart"
            height={this.props.height}
            width={this.props.width}>
          </canvas>
        </div>
      </div>
    );
  }
}

export default ChartWidget;
