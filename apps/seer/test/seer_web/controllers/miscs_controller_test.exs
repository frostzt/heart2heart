defmodule SeerWeb.MiscsControllerTest do
  use SeerWeb.ConnCase

  import Seer.MiscFixtures

  alias Seer.Misc.Miscs

  @create_attrs %{

  }
  @update_attrs %{

  }
  @invalid_attrs %{}

  setup %{conn: conn} do
    {:ok, conn: put_req_header(conn, "accept", "application/json")}
  end

  describe "index" do
    test "lists all misc", %{conn: conn} do
      conn = get(conn, ~p"/api/misc")
      assert json_response(conn, 200)["data"] == []
    end
  end

  describe "create miscs" do
    test "renders miscs when data is valid", %{conn: conn} do
      conn = post(conn, ~p"/api/misc", miscs: @create_attrs)
      assert %{"id" => id} = json_response(conn, 201)["data"]

      conn = get(conn, ~p"/api/misc/#{id}")

      assert %{
               "id" => ^id
             } = json_response(conn, 200)["data"]
    end

    test "renders errors when data is invalid", %{conn: conn} do
      conn = post(conn, ~p"/api/misc", miscs: @invalid_attrs)
      assert json_response(conn, 422)["errors"] != %{}
    end
  end

  describe "update miscs" do
    setup [:create_miscs]

    test "renders miscs when data is valid", %{conn: conn, miscs: %Miscs{id: id} = miscs} do
      conn = put(conn, ~p"/api/misc/#{miscs}", miscs: @update_attrs)
      assert %{"id" => ^id} = json_response(conn, 200)["data"]

      conn = get(conn, ~p"/api/misc/#{id}")

      assert %{
               "id" => ^id
             } = json_response(conn, 200)["data"]
    end

    test "renders errors when data is invalid", %{conn: conn, miscs: miscs} do
      conn = put(conn, ~p"/api/misc/#{miscs}", miscs: @invalid_attrs)
      assert json_response(conn, 422)["errors"] != %{}
    end
  end

  describe "delete miscs" do
    setup [:create_miscs]

    test "deletes chosen miscs", %{conn: conn, miscs: miscs} do
      conn = delete(conn, ~p"/api/misc/#{miscs}")
      assert response(conn, 204)

      assert_error_sent 404, fn ->
        get(conn, ~p"/api/misc/#{miscs}")
      end
    end
  end

  defp create_miscs(_) do
    miscs = miscs_fixture()
    %{miscs: miscs}
  end
end
