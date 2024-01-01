defmodule Seer.Repo do
  use Ecto.Repo,
    otp_app: :seer,
    adapter: Ecto.Adapters.Postgres
end
