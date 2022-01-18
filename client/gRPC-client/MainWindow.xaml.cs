using Grpc.Core;
using Grpc.Net.Client;
using System.Threading;
using System.Windows;

namespace gRPC_client
{
    /// <summary>
    /// Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class MainWindow : Window
    {
        private GrpcChannel _channel;
        private Randomizer.RandomizerClient _client;

        public MainWindow()
        {
            InitializeComponent();
            streamButton.IsEnabled = false;
            requestButton.IsEnabled = false;
        }

        private void connectClick(object sender, RoutedEventArgs e)
        {
            _channel = GrpcChannel.ForAddress($"http://{serverHost.Text}:{serverPort.Text}");
            _client = new Randomizer.RandomizerClient(_channel);

            streamButton.IsEnabled = true;
            requestButton.IsEnabled = true;
            connectButton.IsEnabled = false;
        }

        private async void startStreamClick(object sender, RoutedEventArgs e)
        {
            var cts = new CancellationTokenSource();
            using var valueReplies = _client.GetValueStream(
                new RangeQuery()
                {
                    Count = int.Parse(valueCount.Text),
                    MaxRange = int.Parse(valueMax.Text),
                    MinRange = int.Parse(valueMin.Text)
                },
                cancellationToken: cts.Token
            );

            await foreach (var data in valueReplies.ResponseStream.ReadAllAsync(cts.Token)) {
                outputBox.Text += $"Received value from stream: {data.Value}\n";

            }
        }

        private async void requestClick(object sender, RoutedEventArgs e)
        {
            var reply = await _client.GetValuePairAsync(
                new RangeQuery() {  
                    Count = int.Parse(valueCount.Text),
                    MaxRange = int.Parse(valueMax.Text),
                    MinRange = int.Parse(valueMin.Text)
                }
            );

            outputBox.Text += $"Received {reply.Values.Count} values\n";
            foreach(var value in reply.Values) {
                outputBox.Text += $" => {value.Value}\n";
            }
        }
    }
}
