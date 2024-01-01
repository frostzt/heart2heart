defmodule Seer.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    children = [
      SeerWeb.Telemetry,
      Seer.Repo,
      {DNSCluster, query: Application.get_env(:seer, :dns_cluster_query) || :ignore},
      {Phoenix.PubSub, name: Seer.PubSub},
      # Start the Finch HTTP client for sending emails
      {Finch, name: Seer.Finch},
      # Start a worker by calling: Seer.Worker.start_link(arg)
      # {Seer.Worker, arg},
      # Start to serve requests, typically the last entry
      SeerWeb.Endpoint
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: Seer.Supervisor]
    Supervisor.start_link(children, opts)
  end

  # Tell Phoenix to update the endpoint configuration
  # whenever the application is updated.
  @impl true
  def config_change(changed, _new, removed) do
    SeerWeb.Endpoint.config_change(changed, removed)
    :ok
  end
end
