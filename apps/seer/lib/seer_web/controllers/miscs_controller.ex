defmodule SeerWeb.MiscsController do
  use SeerWeb, :controller

  action_fallback SeerWeb.FallbackController

  def health_check(conn, _params) do
    send_resp(conn, 200, "Seer api up and running...")
  end
end
