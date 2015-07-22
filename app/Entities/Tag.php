<?php

namespace newsletters\Entities;

use Illuminate\Database\Eloquent\Model;
use Prettus\Repository\Contracts\Transformable;
use Prettus\Repository\Traits\TransformableTrait;

class Tag extends Model implements Transformable
{
    use TransformableTrait;

    protected $table = 'tags';

    protected $fillable = [
        'name',
    ];

    public function campaigns()
    {
        return $this->belongsToMany('newsletters\Entities\Campaign', 'campaigns_tags', 'tag_id',
            'campaign_id')->withTimestamps();
    }
}
