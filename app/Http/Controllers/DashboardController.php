<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 15.7.15
 * Time: 18:44
 */

namespace newsletters\Http\Controllers;

class DashboardController extends Controller
{
    public function __construct()
    {
        $this->middleware('auth');
    }

    public function getIndex()
    {
        return view('dashboard.campaigns.list')->with('activeSidebar', 'dashboard');
    }

    public function getNewCampaign()
    {
        return view('dashboard.campaigns.create_new')->with('activeSidebar', 'new-campaign');
    }
}