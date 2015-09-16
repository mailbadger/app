<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 15.7.15
 * Time: 18:44.
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

    public function getTemplates()
    {
        return view('dashboard.templates.list')->with('activeSidebar', 'templates');
    }

    public function getNewTemplate()
    {
        return view('dashboard.templates.create_new')->with('activeSidebar', 'new-template');
    }

    public function getSubscribers()
    {
        return view('dashboard.subscribers.list')->with('activeSidebar', 'sub-lists');
    }

    public function getNewSubscribers()
    {
        return view('dashboard.subscribers.create_new')->with('activeSidebar', 'new-subs');
    }
}