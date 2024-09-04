import sys
import yaml

def generate_compose(output_file, clients):
    """
    Generates a Docker Compose file with a server and a number of clients.
    Args:
        output_file (str): Path to the output file.
        clients (int): Number of clients to generate.

    Returns:
        None

    Example:
        generate_compose("docker-compose-dev.yaml", 5)
    """
    compose = {
        'name': 'tp0',
        'services': {
            'server': {
                'networks': ['testing_net'],
                'image': 'server:latest',
                'container_name': 'server',
                'entrypoint': 'python3 /main.py',
                'environment': [
                    'PYTHONUNBUFFERED=1',
                    'LOGGING_LEVEL=DEBUG',
                    'AMMOUNT_CLIENTS={}'.format(clients),
                ],
                'volumes': [
                    './server/config.ini:/config.ini'
                ],
            }
        },
        'networks': {
            'testing_net': {
                'ipam': {
                    'config': [{'subnet': '172.25.125.0/24'}],
                    'driver': 'default',
                }
            }
        }
    }
    for i in range(1, clients + 1):
        client_name = f'client{i}'
        compose['services'][client_name] = {
            'container_name': client_name,
            'image': 'client:latest',
            'entrypoint': '/client',
            'environment': [
                'PYTHONUNBUFFERED=1',
                'LOGGING_LEVEL=DEBUG',
                'CLI_ID={}'.format(i),
                'FILE_PATH=/agency-{}.csv'.format(i),
            ],
            'networks': ['testing_net'],
            'depends_on': ['server'],
            'volumes': [
                './client/config.yaml:/config.yaml'
            ],
        }

    with open(output_file, 'w') as file:
        yaml.dump(compose, file, default_flow_style=False)

if __name__ == "__main__":
    output_file = sys.argv[1]
    clients = int(sys.argv[2])

    generate_compose(output_file, clients)