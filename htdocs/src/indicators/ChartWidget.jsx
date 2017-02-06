import React from 'react';
import Chart from 'chart.js'

class ChartWidget extends React.Component {
  constructor(props) {
    super(props);
    this.chart = false;
  }

  componentDidMount() {
    if( ! this.chart ) {
      this.chart = new Chart( document.getElementById(this.props.id), {
        type: "pie",
        options: {
          responsive: true,
          responsiveAnimationDuration: 500,
          padding: 10
        }
      });
      this.forceUpdate();
    }
  }

  componentWillUnmount() {
    this.chart.destroy();
    this.chart = null;
  }

  render() {
    // Update existing chart when there's data available
    if( this.chart && Object.keys(this.props.data).length > 0 ) {
      let newData = this.props.reducer(this.props.data);
      this.chart.data.datasets = newData.datasets;
      this.chart.data.labels = newData.labels;
      this.chart.update();
    }

    return (
      <div className="chart-widget">
        {/* Set title if any */}
        {this.props.title &&
          <h2 className="block-title" onClick={() => this.props.onSelection(this.props.reducer)}>
            {this.props.title}
          </h2>
        }

        {/* Set description if any */}
        {this.props.description &&
          <div className="bg-gray">
            <p>{this.props.description}</p>
          </div>
        }

        {/* Chart canvas */}
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
